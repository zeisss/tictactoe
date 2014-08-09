package storage

import (
	"../game"

	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/go-etcd/etcd"
)

func NewEtcdStorage(peer, prefix string, ttl uint64) *etcdStorage {
	client := etcd.NewClient([]string{peer})

	client.SyncCluster()

	return &etcdStorage{client, prefix, ttl}
}

type etcdStorage struct {
	client *etcd.Client
	prefix string
	ttl    uint64
}

func (s *etcdStorage) Get(gameId string) (game.GameState, error) {
	var result game.GameState

	rawResp, err := s.client.RawGet(s.Path(gameId), false, false)
	if err != nil {
		return result, err
	}

	if rawResp.StatusCode == http.StatusNotFound {
		return result, NotFoundError
	}

	resp, err := rawResp.Unmarshal()
	if err != nil {
		return result, err
	}
	if resp.Action != "get" {
		return result, fmt.Errorf("Unexpected response from etcd: action=%s", resp.Action)
	}

	s.Unmarshal(resp.Node.Value, &result)

	return result, nil
}

func (s *etcdStorage) Put(gameId string, game game.GameState) error {
	_, err := s.client.Set(s.Path(gameId), s.Marshal(&game), s.ttl)
	return err
}

func (s *etcdStorage) Path(gameId string) string {
	return fmt.Sprintf("/%s/game/%s", s.prefix, gameId)
}

func (s *etcdStorage) Marshal(game *game.GameState) string {
	data, err := json.Marshal(game)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (s *etcdStorage) Unmarshal(value string, game *game.GameState) {
	if err := json.Unmarshal([]byte(value), game); err != nil {
		panic(err)
	}
}
