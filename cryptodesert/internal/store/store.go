package store

import (
	"crypto-desert/internal/characters"
	"crypto-desert/internal/game"
	"crypto-desert/internal/items"
	"crypto-desert/internal/missions"
	"fmt"
	"sync"
)

// ── CharacterStore ────────────────────────────────────────────────────────────

type CharacterStore struct {
	mu     sync.RWMutex
	chars  map[int]*characters.Character
	nextID int
}

func NewCharacterStore() *CharacterStore {
	return &CharacterStore{chars: make(map[int]*characters.Character), nextID: 1}
}

func (s *CharacterStore) Create(c *characters.Character) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = s.nextID
	s.nextID++
	s.chars[c.ID] = c
}

func (s *CharacterStore) Get(id int) (*characters.Character, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.chars[id]
	if !ok {
		return nil, fmt.Errorf("character %d not found", id)
	}
	return c, nil
}

func (s *CharacterStore) List() []*characters.Character {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*characters.Character, 0, len(s.chars))
	for _, c := range s.chars {
		list = append(list, c)
	}
	return list
}

func (s *CharacterStore) Update(c *characters.Character) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.chars[c.ID]; !ok {
		return fmt.Errorf("character %d not found", c.ID)
	}
	s.chars[c.ID] = c
	return nil
}

func (s *CharacterStore) Delete(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.chars, id)
}

// ── BattleStore ───────────────────────────────────────────────────────────────

type BattleSession struct {
	Battle    *game.Battle
	PlayerID  int
	EnemyName string
}

type BattleStore struct {
	mu       sync.RWMutex
	sessions map[string]*BattleSession
}

func NewBattleStore() *BattleStore {
	return &BattleStore{sessions: make(map[string]*BattleSession)}
}

func (s *BattleStore) Set(id string, sess *BattleSession) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[id] = sess
}

func (s *BattleStore) Get(id string) (*BattleSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[id]
	if !ok {
		return nil, fmt.Errorf("battle session %q not found", id)
	}
	return sess, nil
}

func (s *BattleStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, id)
}

// ── InventoryStore ────────────────────────────────────────────────────────────

type InventoryStore struct {
	mu          sync.RWMutex
	inventories map[int]*items.Inventory
}

func NewInventoryStore() *InventoryStore {
	return &InventoryStore{inventories: make(map[int]*items.Inventory)}
}

func (s *InventoryStore) Get(charID int) *items.Inventory {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.inventories[charID]
}

func (s *InventoryStore) Set(charID int, inv *items.Inventory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.inventories[charID] = inv
}

func (s *InventoryStore) Delete(charID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.inventories, charID)
}

// ── RunnerStore ───────────────────────────────────────────────────────────────

type RunnerStore struct {
	mu      sync.RWMutex
	runners map[string]*missions.Runner
}

func NewRunnerStore() *RunnerStore {
	return &RunnerStore{runners: make(map[string]*missions.Runner)}
}

func (s *RunnerStore) Get(id string) (*missions.Runner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.runners[id]
	if !ok {
		return nil, fmt.Errorf("runner session %q not found", id)
	}
	return r, nil
}

func (s *RunnerStore) Set(id string, r *missions.Runner) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.runners[id] = r
}

func (s *RunnerStore) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.runners, id)
}
