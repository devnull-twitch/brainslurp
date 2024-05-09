package database

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Read(src *badger.Item, target protoreflect.ProtoMessage) error {
	itemBytes, err := src.ValueCopy(nil)
	if err != nil {
		return fmt.Errorf("error copying value bytes: %w", err)
	}

	if err := proto.Unmarshal(itemBytes, target); err != nil {
		return fmt.Errorf("error unmarshaling bytes into proto message: %w", err)
	}

	return nil
}
