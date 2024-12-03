package lws

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/IPA-CyberLab/h132/pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func ReadLWS() (*pb.LetterWritingSet, error) {
	fname := GetLWSWireProtoPath()

	bs, err := os.ReadFile(fname)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", fname, err)
	}

	var lws pb.LetterWritingSet
	if err := proto.Unmarshal(bs, &lws); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file %q: %w", fname, err)
	}

	return &lws, nil
}

func WriteLWS(lws *pb.LetterWritingSet, openflags int) error {
	bs, err := proto.Marshal(lws)
	if err != nil {
		return fmt.Errorf("failed to marshal proto: %w", err)
	}

	f, err := os.OpenFile(GetLWSWireProtoPath(), openflags, 0600)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(bs); err != nil {
		return fmt.Errorf("failed to write proto: %w", err)
	}

	return nil
}

type UpdateFlags int

const (
	UpdateRemoveKey UpdateFlags = 1 << iota
)

func UpdateLWS(lws *pb.LetterWritingSet, flags UpdateFlags) error {
	s := zap.S()

	old, err := ReadLWS()
	if err != nil {
		return err
	}

	// Check for breaking changes
	if lws.Name == "" {
		return fmt.Errorf("name of letter writing set cannot be empty")
	}

	// Check that key are not removed and their names are unique
	oldKeyNameSeen := make(map[string]bool)
	for _, k := range old.Keys {
		oldKeyNameSeen[k.Name] = false
	}
	newKeyNames := make(map[string]struct{})
	for _, k := range lws.Keys {
		if _, ok := newKeyNames[k.Name]; ok {
			return fmt.Errorf("key %q is not unique", k.Name)
		}
		newKeyNames[k.Name] = struct{}{}

		if _, ok := oldKeyNameSeen[k.Name]; ok {
			oldKeyNameSeen[k.Name] = true
		} else {
			s.Infof("UpdateLWS: key %q is newly added to the lws.", k.Name)
		}
	}
	if flags&UpdateRemoveKey == 0 {
		for k, seen := range oldKeyNameSeen {
			if !seen {
				return fmt.Errorf("key %q is removed", k)
			}
		}
	}

	return WriteLWS(lws, os.O_TRUNC|os.O_WRONLY)
}

func GetKeyByName(lws *pb.LetterWritingSet, name string) *pb.KeyImpl {
	for _, k := range lws.Keys {
		if k.Name == name {
			return k
		}
	}
	return nil
}

func GetKeyByPublicKey(lws *pb.LetterWritingSet, needle *ecdsa.PublicKey) *pb.KeyImpl {
	for _, k := range lws.Keys {
		pub := pb.ProtoToPub(k.PublicKey)

		if pub.Equal(needle) {
			return k
		}
	}
	return nil
}

func SetProperty(lws *pb.LetterWritingSet, propertyName, value string) error {
	switch propertyName {
	case "name":
		lws.Name = value
	case "pre_edit_hook":
		lws.PreEditHook = value
	case "post_edit_hook":
		lws.PostEditHook = value
	default:
		return fmt.Errorf("unknown property name: %q", propertyName)
	}
	return nil
}

func DumpStatus(lws *pb.LetterWritingSet) error {
	s := zap.S()

	s.Infof("Letter Writing Set: %s", lws.Name)
	s.Infof("Number of Keys: %d", len(lws.Keys))
	s.Infof("Pre-Edit Hook: %q", lws.PreEditHook)
	s.Infof("Post-Edit Hook: %q", lws.PostEditHook)
	return nil
}
