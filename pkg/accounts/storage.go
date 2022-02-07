package accounts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type (
	Storage interface {
		Add(username string) error
		Get(userID int64) (*AccountObj, error)
		GetByUsername(username string) (*AccountObj, error)
		GetCharacters(userID int64) ([]*GameCharacter, error)
		GetCharacterByName(userID int64, characterName string) (*GameCharacter, error)
		AddCharacter(userID int64, characterName, baseColor string) error
		AddItem(item *InventorySlot) error
		RemoveItem(charID int64, itemID int64, quantity int64) error
		UpdateItem(charID int64, itemID int64, quantity int64) error
		GetItems(charID int64) ([]*InventorySlot, error)
		ChangeCurrentZone(charID int64, newZone string) error
		UnslotEquipment(charID int64, slotID int64, itemID int64) error
		SlotEquipment(charID int64, slotID int64, itemID int64) error
	}

	GameCharacter struct {
		ID          int64
		CurrentZone string
		Name        string
		BaseColor   string
	}

	AccountObj struct {
		ID           int64
		Username     string
		PasswordHash string
	}

	InventorySlot struct {
		ItemID      int64
		CharacterID int64
		Quantity    int64
		SlotID      int64
	}

	databaseStorgae struct {
		conn *pgxpool.Pool
	}
)

func NewStorage(conn *pgxpool.Pool) Storage {
	return &databaseStorgae{
		conn: conn,
	}
}

func (m *databaseStorgae) Add(username string) error {
	_, err := m.conn.Exec(context.Background(), "INSERT INTO accounts ( username ) VALUES ( LOWER($1) )", username)
	return err
}

func (m *databaseStorgae) AddCharacter(userID int64, characterName, baseColor string) error {
	_, err := m.conn.Exec(
		context.Background(),
		`INSERT INTO characters 
			( account_id, character_name, character_display, base_color, current_zone ) VALUES 
			( $1, LOWER($2), $2, $3, $4 )`,
		userID,
		characterName,
		baseColor,
		"starting_zone",
	)
	return err
}

func (m *databaseStorgae) Get(userID int64) (*AccountObj, error) {
	row := m.conn.QueryRow(context.Background(), "SELECT username FROM accounts WHERE id = $1", userID)
	var username string
	if err := row.Scan(&username); err != nil {
		return nil, fmt.Errorf("error selecting account data: %w", err)
	}

	return &AccountObj{
		Username: username,
		ID:       userID,
	}, nil
}

func (m *databaseStorgae) GetByUsername(reqName string) (*AccountObj, error) {
	row := m.conn.QueryRow(context.Background(), "SELECT id, username, password FROM accounts WHERE username = LOWER($1)", reqName)
	var (
		userID       int64
		username     string
		passwordHash *sql.NullString = &sql.NullString{}
	)
	if err := row.Scan(&userID, &username, passwordHash); err != nil {
		return nil, fmt.Errorf("error selecting account data: %w", err)
	}

	return &AccountObj{
		ID:           userID,
		Username:     username,
		PasswordHash: passwordHash.String,
	}, nil
}

func (m *databaseStorgae) GetCharacters(userID int64) ([]*GameCharacter, error) {
	rows, err := m.conn.Query(context.Background(), "SELECT id, character_display, base_color, current_zone FROM characters WHERE account_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error loading characters: %w", err)
	}

	characters := make([]*GameCharacter, 0)
	for rows.Next() {
		var (
			charID        int64
			characterName string
			baseColor     string
			currentZone   string
		)
		if err := rows.Scan(&charID, &characterName, &baseColor, &currentZone); err != nil {
			logrus.WithFields(logrus.Fields{
				"error":      err,
				"account_id": userID,
			}).Error("unable to load chracters")
			continue
		}

		characters = append(characters, &GameCharacter{
			ID:          charID,
			Name:        characterName,
			BaseColor:   baseColor,
			CurrentZone: currentZone,
		})
	}

	return characters, nil
}

func (m *databaseStorgae) AddItem(item *InventorySlot) error {
	_, err := m.conn.Exec(
		context.Background(),
		"INSERT INTO character_items (chracter_id, item_id, quantity) VALUES ( $1, $2, $3 )",
		item.CharacterID,
		item.ItemID,
		item.Quantity,
	)
	if err != nil {
		return fmt.Errorf("error adding item to inventory: %w", err)
	}

	return nil
}

func (m *databaseStorgae) GetItems(charID int64) ([]*InventorySlot, error) {
	rows, err := m.conn.Query(context.Background(), "SELECT item_id, quantity, slot_id FROM character_items WHERE chracter_id = $1", charID)
	if err != nil {
		return nil, fmt.Errorf("error loading items from inventory: %w", err)
	}

	itemIDs := make([]*InventorySlot, 0)
	for rows.Next() {
		var (
			itemID   int64
			quantity int64
			slotID   int64
		)
		if err := rows.Scan(&itemID, &quantity, &slotID); err != nil {
			logrus.WithError(err).Error("unable to load item id")
			continue
		}

		itemIDs = append(itemIDs, &InventorySlot{
			ItemID:      itemID,
			CharacterID: charID,
			Quantity:    quantity,
			SlotID:      slotID,
		})
	}

	return itemIDs, nil
}

func (m *databaseStorgae) GetCharacterByName(userID int64, reqCharacterName string) (*GameCharacter, error) {
	row := m.conn.QueryRow(
		context.Background(),
		"SELECT id, character_display, base_color, current_zone FROM characters WHERE account_id = $1 AND character_name = LOWER($2)",
		userID,
		reqCharacterName,
	)
	var (
		charID        int64
		characterName string
		baseColor     string
		currentZone   string
	)
	if err := row.Scan(&charID, &characterName, &baseColor, &currentZone); err != nil {
		return nil, fmt.Errorf("error loading character: %w", err)
	}

	return &GameCharacter{
		ID:          charID,
		Name:        characterName,
		BaseColor:   baseColor,
		CurrentZone: currentZone,
	}, nil
}

func (m *databaseStorgae) ChangeCurrentZone(charID int64, newZone string) error {
	_, err := m.conn.Exec(context.Background(), "UPDATE characters SET current_zone = $1 WHERE id = $2", newZone, charID)
	if err != nil {
		return fmt.Errorf("error updating character zone: %w", err)
	}

	return nil
}

func (m *databaseStorgae) UnslotEquipment(charID int64, slotID int64, itemID int64) error {
	cmd, err := m.conn.Exec(
		context.Background(),
		"UPDATE character_items SET slot_id = 0 WHERE chracter_id = $1 AND item_id = $2 AND slot_id = $3",
		charID,
		itemID,
		slotID,
	)
	if err != nil {
		return fmt.Errorf("error clearing character inventory slot: %w", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("item was not unequiped")
	}

	return nil
}

func (m *databaseStorgae) SlotEquipment(charID int64, slotID int64, itemID int64) error {
	cmd, err := m.conn.Exec(
		context.Background(),
		`UPDATE character_items SET slot_id = $3 WHERE 
			id = (SELECT id FROM character_items WHERE chracter_id = $1 AND item_id = $2 AND slot_id = 0 LIMIT 1)`,
		charID,
		itemID,
		slotID,
	)
	if err != nil {
		return fmt.Errorf("error updating character inventory slots: %w", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("item was not equiped")
	}

	return nil
}

func (m *databaseStorgae) RemoveItem(charID int64, itemID int64, quantity int64) error {
	cmd, err := m.conn.Exec(
		context.Background(),
		"UPDATE character_items SET quantity = quantity - $3 WHERE id = (SELECT id FROM character_items WHERE chracter_id = $1 AND item_id = $2 LIMIT 1)",
		charID,
		itemID,
		quantity,
	)
	if err != nil {
		return fmt.Errorf("error updating character inventory slots: %w", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("item did not exist")
	}

	_, err = m.conn.Exec(
		context.Background(),
		"DELETE FROM character_items WHERE chracter_id = $1 AND item_id = $2 AND quantity = 0",
		charID,
		itemID,
	)
	if err != nil {
		return fmt.Errorf("error deleting character inventory slots: %w", err)
	}

	return nil
}

func (m *databaseStorgae) UpdateItem(charID int64, itemID int64, quantity int64) error {
	cmd, err := m.conn.Exec(
		context.Background(),
		"UPDATE character_items SET quantity = $3 WHERE id = (SELECT id FROM character_items WHERE chracter_id = $1 AND item_id = $2 LIMIT 1)",
		charID,
		itemID,
		quantity,
	)
	if err != nil {
		return fmt.Errorf("error updating character inventory: %w", err)
	}

	if cmd.RowsAffected() != 1 {
		return fmt.Errorf("item was not updated")
	}

	return nil
}
