package router

import (
	"v2ray.com/core/common/dice"
)

// RandomStrategy represents a random balancing strategy
type RandomStrategy struct {
	balancer *Balancer
}

// PickOutbound implements the BalancingStrategy.
// It picks an outbound from all tags (or alive tags if health check enabled) randomly
func (s *RandomStrategy) PickOutbound() (string, error) {
	tags, err := s.SelectOutbounds()
	if err != nil {
		return "", err
	}
	count := len(tags)
	if count == 0 {
		// goes to fallbackTag
		return "", nil
	}
	return tags[dice.Roll(count)], nil
}

// SelectOutbounds implements BalancingStrategy
func (s *RandomStrategy) SelectOutbounds() ([]string, error) {
	tags, err := s.balancer.SelectOutbounds()
	if err != nil || len(tags) == 0 {
		return nil, err
	}
	if !s.balancer.healthChecker.Settings.Enabled {
		return tags, nil
	}
	aliveTags := make([]string, 0)
	s.balancer.healthChecker.access.Lock()
	defer s.balancer.healthChecker.access.Unlock()
	for _, tag := range tags {
		r, ok := s.balancer.healthChecker.Results[tag]
		if !ok {
			continue
		}
		if r.AverageRTT <= 0 {
			continue
		}
		aliveTags = append(aliveTags, tag)
	}
	if len(aliveTags) == 0 {
		newError("random: no outbound alive").AtInfo().WriteToLog()
	}
	return aliveTags, nil
}
