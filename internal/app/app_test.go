package app

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	ab := New(10, 100, 1000)
	for i := 0; i < 10; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.1"})
		require.NoError(t, err)
	}
	err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.1"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.1"})
	require.NoError(t, err)

	ab.AddToBlackList(context.TODO(), NetworkInfo{Ip: "192.168.1.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", Ip: "192.168.1.1"})
	require.Error(t, err)

	ab.DelFromBlackList(context.TODO(), NetworkInfo{Ip: "192.168.1.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", Ip: "192.168.1.1"})
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User_" + strconv.Itoa(i), Password: "Pass2", Ip: "192.168.1.2"})
		require.NoError(t, err)
	}
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", Ip: "192.168.1.2"})
	require.Error(t, err)

	ab.AddToWhiteList(context.TODO(), NetworkInfo{Ip: "192.168.1.2"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", Ip: "192.168.1.2"})
	require.NoError(t, err)

	ab.DelFromWhiteList(context.TODO(), NetworkInfo{Ip: "192.168.1.2"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", Ip: "192.168.1.2"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "", Password: "Pass2", Ip: "192.168.1.2"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", Ip: "192.168.1.2"})
	require.NoError(t, err)

	for i := 0; i < 1000; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User_" + strconv.Itoa(i), Password: "Pass_" + strconv.Itoa(i), Ip: "192.168.1.3"})
		require.NoError(t, err)
	}
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_1000", Password: "Pass_1000", Ip: "192.168.1.3"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.3"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", Ip: "192.168.1.3"})
	require.NoError(t, err)
}
