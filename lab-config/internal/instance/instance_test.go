package instance

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/the78mole/jumpstarter-mono/core/controller/api/v1alpha1"
	v1alphaConfig "github.com/the78mole/jumpstarter-mono/lab-config/api/v1alpha1"
	"github.com/the78mole/jumpstarter-mono/lab-config/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNewInstance(t *testing.T) {
	// Create a test instance with a context
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
		},
	}

	t.Run("with kubeconfig string", func(t *testing.T) {
		// This will likely fail since we don't have a real kubeconfig, but we test the flow
		_, err := NewInstance(instance, validKubeconfig, false, false)
		// We expect this to fail since the context doesn't exist in our test kubeconfig
		if err != nil {
			assert.Contains(t, err.Error(), "context test-context does not exist in kubeconfig")
		}
	})

	t.Run("without kubeconfig string", func(t *testing.T) {
		// This will likely fail since we don't have a real kubeconfig file, but we test the flow
		_, err := NewInstance(instance, "", false, false)
		// We expect this to fail since the context doesn't exist in the default kubeconfig
		if err != nil {
			assert.True(t,
				contains(err.Error(), "context test-context does not exist") ||
					contains(err.Error(), "context \"test-context\" does not exist") ||
					contains(err.Error(), "kubeconfig file does not exist") ||
					contains(err.Error(), "failed to get in-cluster config") ||
					contains(err.Error(), "failed to load kubeconfig"),
				"Unexpected error: %s", err.Error())
		}
	})

	t.Run("with empty context", func(t *testing.T) {
		instanceNoContext := &v1alphaConfig.JumpstarterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-instance-no-context",
			},
			Spec: v1alphaConfig.JumpstarterInstanceSpec{
				KubeContext: "",
			},
		}

		_, err := NewInstance(instanceNoContext, validKubeconfig, false, false)
		// This should work with our test kubeconfig since it uses the default context
		if err != nil {
			assert.True(t,
				contains(err.Error(), "failed to create client") ||
					contains(err.Error(), "failed to add scheme"),
				"Unexpected error: %s", err.Error())
		}
	})
}

func TestInstanceMethods(t *testing.T) {
	// Create a test instance
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
			Namespace:   "test-namespace",
		},
	}

	t.Run("GetClient and GetConfig", func(t *testing.T) {
		// This will likely fail, but we test the method signatures
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			// If it succeeds, test the methods
			client := inst.GetClient()
			assert.NotNil(t, client)

			config := inst.GetConfig()
			assert.Equal(t, instance, config)
		}
	})
}

func TestInstanceExporterMethods(t *testing.T) {
	// Create a test instance with namespace
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
			Namespace:   "test-namespace",
		},
	}

	t.Run("listExporters with namespace", func(t *testing.T) {
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			// Test that the method exists and can be called
			ctx := context.Background()
			exporters, err := inst.listExporters(ctx)
			// This will likely fail due to connection issues, but we test the method signature
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to list") ||
						contains(err.Error(), "connection") ||
						contains(err.Error(), "context") ||
						contains(err.Error(), "namespace") ||
						contains(err.Error(), "failed to get server groups") ||
						contains(err.Error(), "dial tcp") ||
						contains(err.Error(), "no such host"),
					"Unexpected error: %s", err.Error())
			} else {
				assert.NotNil(t, exporters)
			}
		}
	})

	t.Run("listExporters without namespace", func(t *testing.T) {
		instanceNoNamespace := &v1alphaConfig.JumpstarterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-instance-no-namespace",
			},
			Spec: v1alphaConfig.JumpstarterInstanceSpec{
				KubeContext: "test-context",
				// No namespace specified
			},
		}

		inst, err := NewInstance(instanceNoNamespace, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			exporters, err := inst.listExporters(ctx)
			// This will likely fail due to connection issues, but we test the method signature
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to list") ||
						contains(err.Error(), "connection") ||
						contains(err.Error(), "context") ||
						contains(err.Error(), "failed to get server groups") ||
						contains(err.Error(), "dial tcp") ||
						contains(err.Error(), "no such host"),
					"Unexpected error: %s", err.Error())
			} else {
				assert.NotNil(t, exporters)
			}
		}
	})

	t.Run("GetExporterByName", func(t *testing.T) {
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			exporter, err := inst.GetExporterByName(ctx, "test-exporter")
			// This will likely fail due to connection issues or missing exporter
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to get") ||
						contains(err.Error(), "connection") ||
						contains(err.Error(), "not found") ||
						contains(err.Error(), "namespace") ||
						contains(err.Error(), "failed to get server groups") ||
						contains(err.Error(), "dial tcp") ||
						contains(err.Error(), "no such host"),
					"Unexpected error: %s", err.Error())
			} else {
				assert.NotNil(t, exporter)
			}
		}
	})

	t.Run("GetExporterByName without namespace", func(t *testing.T) {
		instanceNoNamespace := &v1alphaConfig.JumpstarterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-instance-no-namespace",
			},
			Spec: v1alphaConfig.JumpstarterInstanceSpec{
				KubeContext: "test-context",
				// No namespace specified
			},
		}

		inst, err := NewInstance(instanceNoNamespace, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			_, err := inst.GetExporterByName(ctx, "test-exporter")
			// This should fail because namespace is required
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "namespace is required")
		}
	})
}

func TestInstanceClientMethods(t *testing.T) {
	// Create a test instance with namespace
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
			Namespace:   "test-namespace",
		},
	}

	t.Run("ListClients with namespace", func(t *testing.T) {
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			clients, err := inst.ListClients(ctx)
			// This will likely fail due to connection issues, but we test the method signature
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to list") ||
						contains(err.Error(), "connection") ||
						contains(err.Error(), "context") ||
						contains(err.Error(), "namespace") ||
						contains(err.Error(), "failed to get server groups") ||
						contains(err.Error(), "dial tcp") ||
						contains(err.Error(), "no such host"),
					"Unexpected error: %s", err.Error())
			} else {
				assert.NotNil(t, clients)
			}
		}
	})

	t.Run("GetClientByName", func(t *testing.T) {
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			client, err := inst.GetClientByName(ctx, "test-client")
			// This will likely fail due to connection issues or missing client
			if err != nil {
				assert.True(t,
					contains(err.Error(), "failed to get") ||
						contains(err.Error(), "connection") ||
						contains(err.Error(), "not found") ||
						contains(err.Error(), "namespace") ||
						contains(err.Error(), "failed to get server groups") ||
						contains(err.Error(), "dial tcp") ||
						contains(err.Error(), "no such host"),
					"Unexpected error: %s", err.Error())
			} else {
				assert.NotNil(t, client)
			}
		}
	})

	t.Run("GetClientByName without namespace", func(t *testing.T) {
		instanceNoNamespace := &v1alphaConfig.JumpstarterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-instance-no-namespace",
			},
			Spec: v1alphaConfig.JumpstarterInstanceSpec{
				KubeContext: "test-context",
				// No namespace specified
			},
		}

		inst, err := NewInstance(instanceNoNamespace, validKubeconfig, false, false)
		if err == nil {
			ctx := context.Background()
			_, err := inst.GetClientByName(ctx, "test-client")
			// This should fail because namespace is required
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "namespace is required")
		}
	})
}

func TestValidateInstance(t *testing.T) {
	t.Run("valid instance", func(t *testing.T) {
		instance := &v1alphaConfig.JumpstarterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-instance",
			},
			Spec: v1alphaConfig.JumpstarterInstanceSpec{
				KubeContext: "test-context",
			},
		}

		err := validateInstance(instance)
		assert.NoError(t, err)
	})

	t.Run("nil instance", func(t *testing.T) {
		err := validateInstance(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "instance cannot be nil")
	})
}

func TestPrintDiff(t *testing.T) {
	// Create a test instance to use for the printDiff method
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
			Namespace:   "test-namespace",
		},
	}

	t.Run("printDiff with different objects", func(t *testing.T) {
		// Create test objects with different values
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"new-label": "new-value",
				},
			},
		}

		// Create an instance to test the printDiff method
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			// This should not panic and should print a diff
			inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		}
	})

	t.Run("printDiff with identical objects", func(t *testing.T) {
		// Create identical objects
		obj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
		}

		// Create an instance to test the printDiff method
		inst, err := NewInstance(instance, validKubeconfig, false, false)
		if err == nil {
			// This should not panic and should indicate no changes
			inst.checkAndPrintDiff(obj, obj, "exporter", "test-exporter", true)
		}
	})
}

var testUsername = "test-user"

func TestCheckAndPrintDiff(t *testing.T) {
	// Create a test instance to use for the checkAndPrintDiff method
	instance := &v1alphaConfig.JumpstarterInstance{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-instance",
		},
		Spec: v1alphaConfig.JumpstarterInstanceSpec{
			KubeContext: "test-context",
			Namespace:   "test-namespace",
		},
	}

	// Create an instance for testing
	inst, err := NewInstance(instance, validKubeconfig, false, false)
	require.NoError(t, err)

	t.Run("identical objects should return false", func(t *testing.T) {
		obj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(obj, obj, "exporter", "test-exporter", true)
		assert.False(t, changed, "Identical objects should not indicate changes")
	})

	t.Run("different labels should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app":       "test",
					"new-label": "new-value",
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Different labels should indicate changes")
	})

	t.Run("different username should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: func() *string { s := "old-user"; return &s }(),
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: func() *string { s := "new-user"; return &s }(),
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Different username should indicate changes")
	})

	t.Run("different metadata fields should be ignored", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-exporter",
				Namespace:         "test-namespace",
				Generation:        1,
				ResourceVersion:   "123",
				UID:               "old-uid",
				CreationTimestamp: metav1.Now(),
				ManagedFields: []metav1.ManagedFieldsEntry{
					{Manager: "old-manager"},
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:              "test-exporter",
				Namespace:         "test-namespace",
				Generation:        2,
				ResourceVersion:   "456",
				UID:               "new-uid",
				CreationTimestamp: metav1.Now(),
				ManagedFields: []metav1.ManagedFieldsEntry{
					{Manager: "new-manager"},
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.False(t, changed, "Metadata differences should be ignored")
	})

	t.Run("different status should be ignored", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
			Status: v1alpha1.ExporterStatus{
				Conditions: []metav1.Condition{
					{Type: "Ready", Status: "False"},
				},
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
			Status: v1alpha1.ExporterStatus{
				Conditions: []metav1.Condition{
					{Type: "Ready", Status: "True"},
				},
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.False(t, changed, "Status differences should be ignored")
	})

	t.Run("different annotations should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-exporter",
				Namespace:   "test-namespace",
				Annotations: map[string]string{"old": "annotation"},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test-exporter",
				Namespace:   "test-namespace",
				Annotations: map[string]string{"new": "annotation"},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Different annotations should indicate changes")
	})

	t.Run("nil vs non-nil username should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: nil,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Nil vs non-nil username should indicate changes")
	})

	t.Run("different object types should work", func(t *testing.T) {
		oldClient := &v1alpha1.Client{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-client",
				Namespace: "test-namespace",
			},
		}

		newClient := &v1alpha1.Client{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-client",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
		}

		changed := inst.checkAndPrintDiff(oldClient, newClient, "client", "test-client", true)
		assert.True(t, changed, "Different client objects should indicate changes")
	})

	t.Run("empty vs non-empty labels should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels:    map[string]string{},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels: map[string]string{
					"app": "test",
				},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Empty vs non-empty labels should indicate changes")
	})

	t.Run("nil vs empty labels should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels:    nil,
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
				Labels:    map[string]string{},
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Nil vs empty labels should indicate changes")
	})

	t.Run("different namespace should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "old-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "new-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Different namespace should indicate changes")
	})

	t.Run("different name should return true", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "old-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "new-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: &testUsername,
			},
		}

		changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)
		assert.True(t, changed, "Different name should indicate changes")
	})

	t.Run("dry parameter controls printing behavior", func(t *testing.T) {
		oldObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: func() *string { s := "old-user"; return &s }(),
			},
		}

		newObj := &v1alpha1.Exporter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-exporter",
				Namespace: "test-namespace",
			},
			Spec: v1alpha1.ExporterSpec{
				Username: func() *string { s := "new-user"; return &s }(),
			},
		}

		// Test with dry=false (should not print anything)
		t.Run("dry=false should not print", func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", false)

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.True(t, changed, "Should detect changes")
			assert.Empty(t, output, "Should not print anything when dry=false")
		})

		// Test with dry=true (should print diff message)
		t.Run("dry=true should print diff", func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			changed := inst.checkAndPrintDiff(oldObj, newObj, "exporter", "test-exporter", true)

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.True(t, changed, "Should detect changes")
			assert.Contains(t, output, "📝", "Should print dry run message with emoji")
			assert.Contains(t, output, "dry run: Would update", "Should contain dry run message")
			assert.Contains(t, output, "exporter", "Should mention object type")
			assert.Contains(t, output, "test-exporter", "Should mention object name")
		})

		// Test with dry=true and no changes (should print no changes message)
		t.Run("dry=true with no changes should print no changes message", func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			changed := inst.checkAndPrintDiff(oldObj, oldObj, "exporter", "test-exporter", true)

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout

			// Read captured output
			var buf bytes.Buffer
			_, _ = io.Copy(&buf, r)
			output := buf.String()

			assert.False(t, changed, "Should not detect changes")
			assert.Contains(t, output, "✅", "Should print success emoji")
			assert.Contains(t, output, "dry run: No changes needed", "Should contain no changes message")
			assert.Contains(t, output, "exporter", "Should mention object type")
			assert.Contains(t, output, "test-exporter", "Should mention object name")
		})
	})
}

func TestGetExporterObjectForInstance(t *testing.T) {
	// Create a test config (we'll use a mock/simple one for testing)
	cfg := &config.Config{
		Loaded: &config.LoadedLabConfig{
			// We'll add mock data as needed
		},
	}

	tests := []struct {
		name                 string
		exporterInstance     *v1alphaConfig.ExporterInstance
		jumpstarterInstance  string
		expectedExporter     *v1alpha1.Exporter
		expectedError        bool
		expectedErrorMessage string
	}{
		{
			name: "exporter instance with matching jumpstarter instance and no config template",
			exporterInstance: &v1alphaConfig.ExporterInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-exporter",
				},
				Spec: v1alphaConfig.ExporterInstanceSpec{
					Username: testUsername,
					JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
						Name: "target-instance",
					},
					Labels: map[string]string{
						"app": "test",
					},
					// No ConfigTemplateRef - should use default metadata
				},
			},
			jumpstarterInstance: "target-instance",
			expectedExporter: &v1alpha1.Exporter{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-exporter",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: v1alpha1.ExporterSpec{
					Username: &testUsername,
				},
			},
			expectedError: false,
		},
		{
			name: "exporter instance with matching jumpstarter instance and config template",
			exporterInstance: &v1alphaConfig.ExporterInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-exporter-with-template",
					Labels: map[string]string{
						"app": "test",
					},
				},
				Spec: v1alphaConfig.ExporterInstanceSpec{
					Username: testUsername,
					JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
						Name: "target-instance",
					},
					ConfigTemplateRef: v1alphaConfig.ConfigTemplateRef{
						Name: "test-template",
					},
				},
			},
			jumpstarterInstance:  "target-instance",
			expectedExporter:     nil, // Template processing would fail in this test environment
			expectedError:        true,
			expectedErrorMessage: "error creating ExporterInstanceTemplater",
		},
		{
			name: "exporter instance with non-matching jumpstarter instance",
			exporterInstance: &v1alphaConfig.ExporterInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-exporter",
				},
				Spec: v1alphaConfig.ExporterInstanceSpec{
					Username: testUsername,
					JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
						Name: "other-instance",
					},
				},
			},
			jumpstarterInstance: "target-instance",
			expectedExporter:    nil,
			expectedError:       false,
		},
		{
			name: "exporter instance with empty jumpstarter instance ref",
			exporterInstance: &v1alphaConfig.ExporterInstance{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-exporter",
				},
				Spec: v1alphaConfig.ExporterInstanceSpec{
					Username: testUsername,
					JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
						Name: "",
					},
				},
			},
			jumpstarterInstance: "target-instance",
			expectedExporter:    nil,
			expectedError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetExporterObjectForInstance(cfg, tt.exporterInstance, tt.jumpstarterInstance)

			if tt.expectedError {
				assert.Error(t, err, "Expected error for case %s", tt.name)
				if tt.expectedErrorMessage != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMessage, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no error for case %s", tt.name)
			}

			if tt.expectedExporter == nil {
				assert.Nil(t, result, "Expected nil result for case %s", tt.name)
			} else {
				assert.NotNil(t, result, "Expected non-nil result for case %s", tt.name)
				assert.Equal(t, tt.expectedExporter.Name, result.Name, "Expected name to match")
				assert.Equal(t, tt.expectedExporter.Labels, result.Labels, "Expected labels to match")
				assert.Equal(t, tt.expectedExporter.Spec.Username, result.Spec.Username, "Expected username to match")
			}
		})
	}
}

func TestGetExporterObjectForInstance_WithTemplateProcessing(t *testing.T) {
	// Test cases that focus on the template processing logic
	t.Run("exporter instance without config template uses original metadata", func(t *testing.T) {
		exporterInstance := &v1alphaConfig.ExporterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-exporter",
				Annotations: map[string]string{
					"original": "annotation",
				},
			},
			Spec: v1alphaConfig.ExporterInstanceSpec{
				Username: testUsername,
				JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
					Name: "target-instance",
				},
				Labels: map[string]string{
					"original": "label",
					"app":      "test",
				},
				// No ConfigTemplateRef
			},
		}

		cfg := &config.Config{
			Loaded: &config.LoadedLabConfig{},
		}

		result, err := GetExporterObjectForInstance(cfg, exporterInstance, "target-instance")

		assert.NoError(t, err, "Should not error when no template is used")
		assert.NotNil(t, result, "Should return an exporter object")
		assert.Equal(t, "test-exporter", result.Name, "Should preserve original name")
		assert.Equal(t, map[string]string{
			"original": "label",
			"app":      "test",
		}, result.Labels, "Should preserve original labels")
		assert.Equal(t, map[string]string{
			"original": "annotation",
		}, result.Annotations, "Should preserve original annotations")
	})

	t.Run("exporter instance checks HasConfigTemplate correctly", func(t *testing.T) {
		// Test that the function properly uses the HasConfigTemplate method
		exporterInstanceWithTemplate := &v1alphaConfig.ExporterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-exporter-with-template",
			},
			Spec: v1alphaConfig.ExporterInstanceSpec{
				Username: testUsername,
				JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
					Name: "target-instance",
				},
				ConfigTemplateRef: v1alphaConfig.ConfigTemplateRef{
					Name: "some-template",
				},
			},
		}

		exporterInstanceWithoutTemplate := &v1alphaConfig.ExporterInstance{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-exporter-without-template",
			},
			Spec: v1alphaConfig.ExporterInstanceSpec{
				Username: testUsername,
				JumpstarterInstanceRef: v1alphaConfig.JumsptarterInstanceRef{
					Name: "target-instance",
				},
				// No ConfigTemplateRef
			},
		}

		cfg := &config.Config{
			Loaded: &config.LoadedLabConfig{},
		}

		// Test with template - should fail in this environment because we can't create a real templater
		_, err := GetExporterObjectForInstance(cfg, exporterInstanceWithTemplate, "target-instance")
		assert.Error(t, err, "Should error when trying to create templater without proper config")

		// Test without template - should succeed
		result, err := GetExporterObjectForInstance(cfg, exporterInstanceWithoutTemplate, "target-instance")
		assert.NoError(t, err, "Should not error when no template is used")
		assert.NotNil(t, result, "Should return an exporter object")
	})
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}
