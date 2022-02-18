package app

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	ab := New(context.TODO(), 10, 100, 1000)
	for i := 0; i < 10; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
		require.NoError(t, err)
	}
	err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	require.NoError(t, err)

	ab.AddToBlackList(context.TODO(), NetworkInfo{IP: "192.168.1.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", IP: "192.168.1.1"})
	require.Error(t, err)

	ab.DelFromBlackList(context.TODO(), NetworkInfo{IP: "192.168.1.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", IP: "192.168.1.1"})
	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User_" + strconv.Itoa(i), Password: "Pass2", IP: "192.168.2.1"})
		require.NoError(t, err)
	}
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.Error(t, err)

	ab.AddToWhiteList(context.TODO(), NetworkInfo{IP: "192.168.2.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.NoError(t, err)

	ab.DelFromWhiteList(context.TODO(), NetworkInfo{IP: "192.168.2.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "", Password: "Pass2", IP: "192.168.2.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.NoError(t, err)

	for i := 0; i < 1000; i++ {
		err := ab.Login(context.TODO(), LoginInfo{
			Login:    "User_" + strconv.Itoa(i),
			Password: "Pass_" + strconv.Itoa(i), IP: "192.168.1.3",
		})
		require.NoError(t, err)
	}
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_1000", Password: "Pass_1000", IP: "192.168.1.3"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.3"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.3"})
	require.NoError(t, err)
}

func TestReset(t *testing.T) {
	ab := New(context.TODO(), 10, 100, 1000)
	for i := 0; i < 10; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
		require.NoError(t, err)
	}
	err := ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.1"})
	require.NoError(t, err)
}

func TestBlackList(t *testing.T) {
	ab := New(context.TODO(), 10, 100, 1000)

	ab.AddToBlackList(context.TODO(), NetworkInfo{IP: "192.168.1.0/25"})
	err := ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", IP: "192.168.1.1"})
	require.Error(t, err)

	ab.DelFromBlackList(context.TODO(), NetworkInfo{IP: "192.168.1.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User2", Password: "Pass1", IP: "192.168.1.1"})
	require.NoError(t, err)
}

func TestWhiteList(t *testing.T) {
	ab := New(context.TODO(), 10, 100, 1000)

	for i := 0; i < 100; i++ {
		err := ab.Login(context.TODO(), LoginInfo{Login: "User_" + strconv.Itoa(i), Password: "Pass2", IP: "192.168.2.1"})
		require.NoError(t, err)
	}
	err := ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.Error(t, err)

	ab.AddToWhiteList(context.TODO(), NetworkInfo{IP: "192.168.2.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.NoError(t, err)

	ab.DelFromWhiteList(context.TODO(), NetworkInfo{IP: "192.168.2.0/25"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User_100", Password: "Pass2", IP: "192.168.2.1"})
	require.Error(t, err)
}

func TestIP(t *testing.T) {
	ab := New(context.TODO(), 10, 100, 1000)

	for i := 0; i < 1000; i++ {
		err := ab.Login(context.TODO(), LoginInfo{
			Login:    "User_" + strconv.Itoa(i),
			Password: "Pass_" + strconv.Itoa(i), IP: "192.168.1.3",
		})
		require.NoError(t, err)
	}
	err := ab.Login(context.TODO(), LoginInfo{Login: "User_1000", Password: "Pass_1000", IP: "192.168.1.3"})
	require.Error(t, err)

	ab.Reset(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.3"})
	err = ab.Login(context.TODO(), LoginInfo{Login: "User1", Password: "Pass1", IP: "192.168.1.3"})
	require.NoError(t, err)
}
