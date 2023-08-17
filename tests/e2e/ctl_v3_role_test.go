package e2e

import (
	"fmt"
	"testing"
)

func TestCtlV3RoleAdd(t *testing.T)          { testCtl(t, roleAddTest) }
func TestCtlV3RoleAddNoTLS(t *testing.T)     { testCtl(t, roleAddTest, withCfg(configNoTLS)) }
func TestCtlV3RoleAddClientTLS(t *testing.T) { testCtl(t, roleAddTest, withCfg(configClientTLS)) }
func TestCtlV3RoleAddPeerTLS(t *testing.T)   { testCtl(t, roleAddTest, withCfg(configPeerTLS)) }
func TestCtlV3RoleAddTimeout(t *testing.T)   { testCtl(t, roleAddTest, withDialTimeout(0)) }

func TestCtlV3RoleGrant(t *testing.T) { testCtl(t, roleGrantTest) }

func roleAddTest(cx ctlCtx) {
	cmdSet := []struct {
		args        []string
		expectedStr string
	}{
		// Add a role.
		{
			args:        []string{"add", "root"},
			expectedStr: "Role root created",
		},
		// Try adding the same role.
		{
			args:        []string{"add", "root"},
			expectedStr: "role name already exists",
		},
	}

	for i, cmd := range cmdSet {
		if err := ctlV3Role(cx, cmd.args, cmd.expectedStr); err != nil {
			if cx.dialTimeout > 0 && !isGRPCTimedout(err) {
				cx.t.Fatalf("roleAddTest #%d: ctlV3Role error (%v)", i, err)
			}
		}
	}
}

func roleGrantTest(cx ctlCtx) {
	cmdSet := []struct {
		args        []string
		expectedStr string
	}{
		// Add a role.
		{
			args:        []string{"add", "root"},
			expectedStr: "Role root created",
		},
		// Grant read permission to the role.
		{
			args:        []string{"grant", "root", "read", "foo"},
			expectedStr: "Role root updated",
		},
		// Grant write permission to the role.
		{
			args:        []string{"grant", "root", "write", "foo"},
			expectedStr: "Role root updated",
		},
		// Grant rw permission to the role.
		{
			args:        []string{"grant", "root", "readwrite", "foo"},
			expectedStr: "Role root updated",
		},
		// Try granting invalid permission to the role.
		{
			args:        []string{"grant", "root", "123", "foo"},
			expectedStr: "invalid permission type",
		},
	}

	for i, cmd := range cmdSet {
		if err := ctlV3Role(cx, cmd.args, cmd.expectedStr); err != nil {
			cx.t.Fatalf("roleGrantTest #%d: ctlV3Role error (%v)", i, err)
		}
	}
}

func ctlV3Role(cx ctlCtx, args []string, expStr string) error {
	cmdArgs := append(cx.PrefixArgs(), "role")
	cmdArgs = append(cmdArgs, args...)

	return spawnWithExpect(cmdArgs, expStr)
}

func ctlV3RoleGrantPermission(cx ctlCtx, rolename string, perm grantingPerm) error {
	cmdArgs := append(cx.PrefixArgs(), "role", "grant-permission")
	if perm.prefix {
		cmdArgs = append(cmdArgs, "--prefix")
	} else if len(perm.rangeEnd) == 1 && perm.rangeEnd[0] == '\x00' {
		cmdArgs = append(cmdArgs, "--from-key")
	}

	cmdArgs = append(cmdArgs, rolename)
	cmdArgs = append(cmdArgs, grantingPermToArgs(perm)...)

	proc, err := spawnCmd(cmdArgs)
	if err != nil {
		return err
	}

	expStr := fmt.Sprintf("Role %s updated", rolename)
	_, err = proc.Expect(expStr)
	return err
}

func ctlV3RoleRevokePermission(cx ctlCtx, rolename string, key, rangeEnd string, fromKey bool) error {
	cmdArgs := append(cx.PrefixArgs(), "role", "revoke-permission")
	cmdArgs = append(cmdArgs, rolename)
	cmdArgs = append(cmdArgs, key)
	var expStr string
	if len(rangeEnd) != 0 {
		cmdArgs = append(cmdArgs, rangeEnd)
		expStr = fmt.Sprintf("Permission of range [%s, %s) is revoked from role %s", key, rangeEnd, rolename)
	} else if fromKey {
		cmdArgs = append(cmdArgs, "--from-key")
		expStr = fmt.Sprintf("Permission of range [%s, <open ended> is revoked from role %s", key, rolename)
	} else {
		expStr = fmt.Sprintf("Permission of key %s is revoked from role %s", key, rolename)
	}

	proc, err := spawnCmd(cmdArgs)
	if err != nil {
		return err
	}

	_, err = proc.Expect(expStr)
	return err
}

type grantingPerm struct {
	read     bool
	write    bool
	key      string
	rangeEnd string
	prefix   bool
}

func grantingPermToArgs(perm grantingPerm) []string {
	permstr := ""

	if perm.read {
		permstr += "read"
	}

	if perm.write {
		permstr += "write"
	}

	if len(permstr) == 0 {
		panic("invalid granting permission")
	}

	if len(perm.rangeEnd) == 0 {
		return []string{permstr, perm.key}
	}

	if len(perm.rangeEnd) == 1 && perm.rangeEnd[0] == '\x00' {
		return []string{permstr, perm.key}
	}

	return []string{permstr, perm.key, perm.rangeEnd}
}
