package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/pachyderm/pachyderm/src/client/pkg/require"
	tu "github.com/pachyderm/pachyderm/src/server/pkg/testutil"
)

var stdoutMutex = &sync.Mutex{}

func TestMetricsNormalDeployment(t *testing.T) {
	// Run deploy normally, should see METRICS=true
	testDeploy(t, false, false, true)
}

func TestMetricsNormalDeploymentNoMetricsFlagSet(t *testing.T) {
	// Run deploy normally, should see METRICS=true
	testDeploy(t, false, true, false)
}

func TestMetricsDevDeployment(t *testing.T) {
	// Run deploy w dev flag, should see METRICS=false
	testDeploy(t, true, false, false)
}

func TestMetricsDevDeploymentNoMetricsFlagSet(t *testing.T) {
	// Run deploy w dev flag, should see METRICS=false
	testDeploy(t, true, true, false)
}

func TestPortForwardError(t *testing.T) {
	os.Setenv("PACHD_ADDRESS", "localhost:30650")
	c := tu.Cmd("pachctl", "version", "--timeout=1ns", "--no-port-forwarding")
	var errMsg bytes.Buffer
	c.Stdout = ioutil.Discard
	c.Stderr = &errMsg
	err := c.Run()
	require.YesError(t, err) // 1ns should prevent even local connections
	require.Matches(t, "port-forward", errMsg.String())
}

func TestNoPortError(t *testing.T) {
	os.Setenv("PACHD_ADDRESS", "127.127.127.0")
	c := tu.Cmd("pachctl", "version", "--timeout=1ns", "--no-port-forwarding")
	var errMsg bytes.Buffer
	c.Stdout = ioutil.Discard
	c.Stderr = &errMsg
	err := c.Run()
	require.YesError(t, err) // 1ns should prevent even local connections
	require.Matches(t, "does not seem to be host:port", errMsg.String())
}

func TestWeirdPortError(t *testing.T) {
	os.Setenv("PACHD_ADDRESS", "localhost:30560")
	c := tu.Cmd("pachctl", "version", "--timeout=1ns", "--no-port-forwarding")
	var errMsg bytes.Buffer
	c.Stdout = ioutil.Discard
	c.Stderr = &errMsg
	err := c.Run()
	require.YesError(t, err) // 1ns should prevent even local connections
	require.Matches(t, "30650", errMsg.String())
}

func testDeploy(t *testing.T, devFlag bool, noMetrics bool, expectedEnvValue bool) {
	//t.Parallel()
	//stdoutMutex.Lock()
	//defer stdoutMutex.Unlock()

	//// Setup user config prior to test
	//// So that stdout only contains JSON no warnings
	//_, err := config.Read()
	//require.NoError(t, err)

	//old := os.Stdout
	//r, w, _ := os.Pipe()
	//os.Stdout = w

	//os.Args = []string{
	//"deploy",
	//"local",
	//"--dry-run",
	//fmt.Sprintf("-d=%v", devFlag),
	//}
	//err = deploycmds.DeployCmd(&noMetrics).Execute()
	//require.NoError(t, err)
	//require.NoError(t, w.Close())
	//// restore stdout
	//os.Stdout = old

	//decoder := json.NewDecoder(r)
	//foundPachdManifest := false
	//// Loop through generated manifest until we find a
	//// ReplicationController (limit of 100 makes sure test
	//// fails quickly if there is no RC)
	//for i := 0; i < 100; i++ {
	//var manifest *extensions.Deployment
	//err = decoder.Decode(&manifest)
	//if err == io.EOF {
	//break
	//}
	//if err != nil {
	//// Not a replication controller
	//continue
	//}
	//require.NoError(t, err)
	//if manifest.ObjectMeta.Name == "pachd" && manifest.Kind == "Deployment" {
	//foundPachdManifest = true
	//expectedMetricEnvVar := api.EnvVar{
	//Name:  "METRICS",
	//Value: fmt.Sprintf("%v", expectedEnvValue),
	//}
	//var env []interface{}
	//require.Equal(t, 1, len(manifest.Spec.Template.Spec.Containers))
	//for _, value := range manifest.Spec.Template.Spec.Containers[0].Env {
	//env = append(env, value)
	//}
	//require.OneOfEquals(t, interface{}(expectedMetricEnvVar), env)
	//}
	//}
	//require.Equal(t, true, foundPachdManifest)
}
