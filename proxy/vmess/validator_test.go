package vmess_test

import (
	"testing"
	"time"

	"github.com/v2fly/v2ray-core/v4/common"
	"github.com/v2fly/v2ray-core/v4/common/protocol"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/common/uuid"
	. "github.com/v2fly/v2ray-core/v4/proxy/vmess"
)

func toAccount(a *Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestUserValidator(t *testing.T) {
	hasher := protocol.DefaultIDHash
	v := NewTimedUserValidator(hasher)
	defer common.Close(v)

	id := uuid.New()
	user := &protocol.MemoryUser{
		Email: "test",
		Account: toAccount(&Account{
			Id: id.String(),
		}),
	}
	common.Must(v.Add(user))

	{
		testSmallLag := func(lag time.Duration) {
			ts := protocol.Timestamp(time.Now().Add(time.Second * lag).Unix())
			idHash := hasher(id.Bytes())
			common.Must2(serial.WriteUint64(idHash, uint64(ts)))
			userHash := idHash.Sum(nil)

			euser, ets, found, _ := v.Get(userHash)
			if !found {
				t.Fatal("user not found")
			}
			if euser.Email != user.Email {
				t.Error("unexpected user email: ", euser.Email, " want ", user.Email)
			}
			if ets != ts {
				t.Error("unexpected timestamp: ", ets, " want ", ts)
			}
		}

		testSmallLag(0)
		testSmallLag(40)
		testSmallLag(-40)
		testSmallLag(80)
		testSmallLag(-80)
		testSmallLag(120)
		testSmallLag(-120)
	}

	{
		testBigLag := func(lag time.Duration) {
			ts := protocol.Timestamp(time.Now().Add(time.Second * lag).Unix())
			idHash := hasher(id.Bytes())
			common.Must2(serial.WriteUint64(idHash, uint64(ts)))
			userHash := idHash.Sum(nil)

			euser, _, found, _ := v.Get(userHash)
			if found || euser != nil {
				t.Error("unexpected user")
			}
		}

		testBigLag(121)
		testBigLag(-121)
		testBigLag(310)
		testBigLag(-310)
		testBigLag(500)
		testBigLag(-500)
	}

	if v := v.Remove(user.Email); !v {
		t.Error("unable to remove user")
	}
	if v := v.Remove(user.Email); v {
		t.Error("remove user twice")
	}
}

func BenchmarkUserValidator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hasher := protocol.DefaultIDHash
		v := NewTimedUserValidator(hasher)

		for j := 0; j < 1500; j++ {
			id := uuid.New()
			v.Add(&protocol.MemoryUser{
				Email: "test",
				Account: toAccount(&Account{
					Id: id.String(),
				}),
			})
		}

		common.Close(v)
	}
}
