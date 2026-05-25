package items

import (
	"fmt"
)

// ── Inventory Entry ───────────────────────────────────────────────────────────

// InventoryEntry representa uma pilha de um item específico no inventário
type InventoryEntry struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

// Item retorna a definição do item desta entrada
func (e InventoryEntry) Item() (Item, bool) {
	return ItemByID(e.ItemID)
}

// ── Equipment Slots ───────────────────────────────────────────────────────────

// EquippedItems representa os 3 slots de equipamento do personagem
type EquippedItems struct {
	Weapon    string `json:"weapon"`
	Armor     string `json:"armor"`
	Accessory string `json:"accessory"`
}

// SlotFor retorna o ItemID equipado em um slot específico
func (e *EquippedItems) SlotFor(slot EquipSlot) string {
	switch slot {
	case SlotWeapon:    return e.Weapon
	case SlotArmor:     return e.Armor
	case SlotAccessory: return e.Accessory
	}
	return ""
}

// SetSlot equipa um item em um slot
func (e *EquippedItems) SetSlot(slot EquipSlot, itemID string) {
	switch slot {
	case SlotWeapon:    e.Weapon = itemID
	case SlotArmor:     e.Armor = itemID
	case SlotAccessory: e.Accessory = itemID
	}
}

// AllEquipped retorna todos os ItemIDs equipados (sem vazios)
func (e *EquippedItems) AllEquipped() []string {
	var result []string
	if e.Weapon != ""    { result = append(result, e.Weapon) }
	if e.Armor != ""     { result = append(result, e.Armor) }
	if e.Accessory != "" { result = append(result, e.Accessory) }
	return result
}

// ── Inventory ─────────────────────────────────────────────────────────────────

// Inventory gerencia os itens e equipamentos de um personagem.
// Separado do Character para manter o pacote characters sem dependência de items.
type Inventory struct {
	CharacterID int              `json:"character_id"`
	Slots       []InventoryEntry `json:"slots"`
	Equipped    EquippedItems    `json:"equipped"`
	MaxSlots    int              `json:"max_slots"`
}

// NewInventory cria um inventário vazio para um personagem
func NewInventory(characterID int) *Inventory {
	return &Inventory{
		CharacterID: characterID,
		Slots:       make([]InventoryEntry, 0),
		MaxSlots:    20,
	}
}

// ── Consultas ─────────────────────────────────────────────────────────────────

// Count retorna a quantidade de um item no inventário
func (inv *Inventory) Count(itemID string) int {
	for _, e := range inv.Slots {
		if e.ItemID == itemID {
			return e.Quantity
		}
	}
	return 0
}

// HasItem retorna true se o inventário tem ao menos 1 do item
func (inv *Inventory) HasItem(itemID string) bool {
	return inv.Count(itemID) > 0
}

// IsEquipped retorna true se o item está atualmente equipado
func (inv *Inventory) IsEquipped(itemID string) bool {
	eq := inv.Equipped
	return eq.Weapon == itemID || eq.Armor == itemID || eq.Accessory == itemID
}

// UsedSlots retorna quantos slots distintos estão ocupados
func (inv *Inventory) UsedSlots() int {
	return len(inv.Slots)
}

// IsFull retorna true se não há mais espaço no inventário
func (inv *Inventory) IsFull() bool {
	return len(inv.Slots) >= inv.MaxSlots
}

// ── Modificações ──────────────────────────────────────────────────────────────

// Add adiciona qty unidades de um item ao inventário.
// Retorna erro se o inventário estiver cheio ou o item não existir.
func (inv *Inventory) Add(itemID string, qty int) error {
	item, ok := ItemByID(itemID)
	if !ok {
		return fmt.Errorf("item %q não encontrado no catálogo", itemID)
	}
	if qty <= 0 {
		return fmt.Errorf("quantidade deve ser > 0")
	}

	// Tenta empilhar em entrada existente
	if item.Stackable {
		for i, e := range inv.Slots {
			if e.ItemID == itemID {
				if item.MaxStack > 0 && inv.Slots[i].Quantity+qty > item.MaxStack {
					return fmt.Errorf("stack máximo (%d) atingido para %s", item.MaxStack, item.Name)
				}
				inv.Slots[i].Quantity += qty
				return nil
			}
		}
	}

	// Novo slot
	if inv.IsFull() {
		return fmt.Errorf("inventário cheio (%d/%d slots)", len(inv.Slots), inv.MaxSlots)
	}
	inv.Slots = append(inv.Slots, InventoryEntry{ItemID: itemID, Quantity: qty})
	return nil
}

// Remove remove qty unidades de um item.
// Retorna erro se não houver quantidade suficiente.
func (inv *Inventory) Remove(itemID string, qty int) error {
	for i, e := range inv.Slots {
		if e.ItemID == itemID {
			if e.Quantity < qty {
				return fmt.Errorf("quantidade insuficiente: tem %d, precisa de %d", e.Quantity, qty)
			}
			inv.Slots[i].Quantity -= qty
			if inv.Slots[i].Quantity == 0 {
				inv.Slots = append(inv.Slots[:i], inv.Slots[i+1:]...)
			}
			return nil
		}
	}
	return fmt.Errorf("item %q não encontrado no inventário", itemID)
}

// ── Equipamento ───────────────────────────────────────────────────────────────

// EquipResult descreve o que mudou ao equipar um item
type EquipResult struct {
	Equipped   Item
	Unequipped *Item // item que foi desequipado (se havia algo no slot)
}

// Equip equipa um item do inventário.
// Retorna o que foi equipado e o que foi desequipado (se havia).
// O item equipado sai do inventário; o desequipado volta para o inventário.
func (inv *Inventory) Equip(itemID string, class string) (EquipResult, error) {
	item, ok := ItemByID(itemID)
	if !ok {
		return EquipResult{}, fmt.Errorf("item %q não encontrado", itemID)
	}
	if item.Slot == SlotNone {
		return EquipResult{}, fmt.Errorf("%s não é um item equipável", item.Name)
	}
	if !item.Classes.Allows(class) {
		return EquipResult{}, fmt.Errorf(
			"%s não pode ser equipado por %s (requer: %v)", item.Name, class, []string(item.Classes),
		)
	}
	if !inv.HasItem(itemID) {
		return EquipResult{}, fmt.Errorf("%s não está no inventário", item.Name)
	}

	result := EquipResult{Equipped: item}

	// Desequipa o que estava no slot (se houver) e devolve ao inventário
	currentID := inv.Equipped.SlotFor(item.Slot)
	if currentID != "" {
		current, _ := ItemByID(currentID)
		result.Unequipped = &current
		// Devolve ao inventário (ignora erro de stack cheio — sempre tem espaço pois removemos um)
		_ = inv.Add(currentID, 1)
	}

	// Remove do inventário e equipa
	if err := inv.Remove(itemID, 1); err != nil {
		return EquipResult{}, err
	}
	inv.Equipped.SetSlot(item.Slot, itemID)

	return result, nil
}

// Unequip remove um item do slot e o devolve ao inventário.
func (inv *Inventory) Unequip(slot EquipSlot) error {
	currentID := inv.Equipped.SlotFor(slot)
	if currentID == "" {
		return fmt.Errorf("slot %s está vazio", slot)
	}

	if inv.IsFull() {
		return fmt.Errorf("inventário cheio — faça espaço antes de desequipar")
	}

	inv.Equipped.SetSlot(slot, "")
	return inv.Add(currentID, 1)
}

// ── Uso de Itens ──────────────────────────────────────────────────────────────

// CharacterStats é a interface mínima que o inventário precisa para aplicar efeitos.
// Implementada pelo *characters.Character — evita import circular.
type CharacterStats interface {
	GetHP() int
	GetMaxHP() int
	GetMana() int
	GetMaxMana() int
	ApplyHeal(amount int)
	ApplyManaRestore(amount int)
	ApplyStatus(effect string, duration, power int)
	IsAlive() bool
}

// UseItem usa um consumível do inventário.
// inBattle indica se estamos dentro de uma batalha.
func (inv *Inventory) UseItem(itemID string, target CharacterStats, inBattle bool) (UseResult, error) {
	item, ok := ItemByID(itemID)
	if !ok {
		return UseResult{}, fmt.Errorf("item %q não encontrado", itemID)
	}
	if !item.Usable {
		return UseResult{}, fmt.Errorf("%s não pode ser usado diretamente", item.Name)
	}
	if inBattle && !item.UsableInBattle {
		return UseResult{}, fmt.Errorf("%s não pode ser usado em batalha", item.Name)
	}
	if !inv.HasItem(itemID) {
		return UseResult{}, fmt.Errorf("%s não está no inventário", item.Name)
	}

	result := UseResult{Success: true}
	ef := item.Effect

	// HP
	if ef.HealHP > 0 {
		result.HPRestored = ef.HealHP
		target.ApplyHeal(ef.HealHP)
	}
	if ef.HealHPPct > 0 {
		amount := int(float64(target.GetMaxHP()) * ef.HealHPPct)
		result.HPRestored = amount
		target.ApplyHeal(amount)
	}

	// Mana
	if ef.HealMana > 0 {
		result.ManaRestored = ef.HealMana
		target.ApplyManaRestore(ef.HealMana)
	}
	if ef.HealManaPct > 0 {
		amount := int(float64(target.GetMaxMana()) * ef.HealManaPct)
		result.ManaRestored = amount
		target.ApplyManaRestore(amount)
	}

	// Status buff
	if ef.AppliesStatus != "" {
		target.ApplyStatus(ef.AppliesStatus, ef.StatusDuration, ef.StatusPower)
	}

	// Consome o item
	if err := inv.Remove(itemID, 1); err != nil {
		return UseResult{}, err
	}

	result.Message = fmt.Sprintf("Usou %s: %s", item.Name, item.Description)
	return result, nil
}

// ── Bônus Totais de Equipamento ───────────────────────────────────────────────

// EquipmentBonuses agrega todos os bônus de stat dos itens atualmente equipados.
// Chamado pelo sistema de combate para calcular stats efetivos.
type EquipmentBonuses struct {
	AttackMod         int
	StrengthMod       int
	CA                int
	Defense           int
	Speed             int
	MaxHP             int
	MaxMana           int
	CryptoFactorBonus float64
}

// TotalBonuses calcula os bônus somados de todos os equipamentos
func (inv *Inventory) TotalBonuses() EquipmentBonuses {
	var b EquipmentBonuses
	for _, id := range inv.Equipped.AllEquipped() {
		item, ok := ItemByID(id)
		if !ok {
			continue
		}
		ef := item.Effect
		b.AttackMod         += ef.BonusAttackMod
		b.StrengthMod       += ef.BonusStrengthMod
		b.CA                += ef.BonusCA
		b.Defense           += ef.BonusDefense
		b.Speed             += ef.BonusSpeed
		b.MaxHP             += ef.BonusMaxHP
		b.MaxMana           += ef.BonusMaxMana
		b.CryptoFactorBonus += ef.CryptoFactorBonus
	}
	return b
}
