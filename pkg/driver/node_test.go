/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/mock/gomock"
	"github.com/kubernetes-sigs/aws-efs-csi-driver/pkg/driver/mocks"
)

const (
	volumeId   = "fs-abc123"
	targetPath = "/target/path"
)

type errtyp struct {
	code    string
	message string
}

func setup(mockCtrl *gomock.Controller) (*mocks.MockMounter, *Driver, context.Context) {
	mockMounter := mocks.NewMockMounter(mockCtrl)
	driver := &Driver{
		endpoint: "endpoint",
		nodeID:   "nodeID",
		mounter:  mockMounter,
	}
	ctx := context.Background()
	return mockMounter, driver, ctx
}

func testResult(t *testing.T, funcName string, ret interface{}, err error, expectError errtyp) {
	if expectError.code == "" {
		if err != nil {
			t.Fatalf("%s is failed: %v", funcName, err)
		}
		if ret == nil {
			t.Fatal("Expected non-nil return value")
		}
	} else {
		if err == nil {
			t.Fatalf("%s is not failed", funcName)
		}
		// Sure would be nice if grpc.statusError was exported :(
		// The error string looks like:
		// "rpc error: code = {code} desc = {desc}"
		tokens := strings.SplitN(err.Error(), " = ", 3)
		expCode := strings.Split(tokens[1], " ")[0]
		if expCode != expectError.code {
			t.Fatalf("Expected error code %q but got %q", expCode, expectError.code)
		}
		if tokens[2] != expectError.message {
			t.Fatalf("\nExpected error message: %s\nActual error message:   %s", expectError.message, tokens[2])
		}
	}
}

func TestNodePublishVolume(t *testing.T) {

	var (
		accessPointID = "fsap-abcd1234"
		stdVolCap     = &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			},
		}
	)

	testCases := []struct {
		name          string
		req           *csi.NodePublishVolumeRequest
		expectMakeDir bool
		mountArgs     []interface{}
		mountSuccess  bool
		expectError   errtyp
	}{
		{
			name: "success: normal",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: empty path",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + ":",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: empty path and access point",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + "::",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: normal with read only mount",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
				Readonly:         true,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", []string{"ro"}},
			mountSuccess:  true,
		},
		{
			name: "success: normal with tls mount options",
			req: &csi.NodePublishVolumeRequest{
				VolumeId: volumeId,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{
							MountFlags: []string{"tls"},
						},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
				TargetPath: targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", []string{"tls"}},
			mountSuccess:  true,
		},
		{
			// TODO: Validate deprecation warning
			name: "success: normal with path in volume context",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
				VolumeContext:    map[string]string{"path": "/a/b"},
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/a/b", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "fail: path in volume context must be absolute",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
				VolumeContext:    map[string]string{"path": "a/b"},
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: `Volume context property "path" must be an absolute path`,
			},
		},
		{
			name: "success: normal with path in volume handle",
			req: &csi.NodePublishVolumeRequest{
				// This also shows that the path gets cleaned
				VolumeId:         volumeId + ":/a/b/",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/a/b", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: normal with path in volume handle, empty access point",
			req: &csi.NodePublishVolumeRequest{
				// This also shows that relative paths are allowed when specified via volume handle
				VolumeId:         volumeId + ":a/b/:",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":a/b", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: path in volume handle takes precedence",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + ":/a/b/",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
				VolumeContext:    map[string]string{"path": "/c/d"},
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/a/b", targetPath, "efs", gomock.Any()},
			mountSuccess:  true,
		},
		{
			name: "success: access point in volume handle, no path",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + "::" + accessPointID,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", []string{"accesspoint=" + accessPointID, "tls"}},
			mountSuccess:  true,
		},
		{
			name: "success: path and access point in volume handle",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + ":/a/b:" + accessPointID,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/a/b", targetPath, "efs", []string{"accesspoint=" + accessPointID, "tls"}},
			mountSuccess:  true,
		},
		{
			// TODO: Validate deprecation warning
			name: "success: same access point in volume handle and mount options",
			req: &csi.NodePublishVolumeRequest{
				VolumeId: volumeId + "::" + accessPointID,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{
							// This also shows we allow the `tls` option to exist already
							MountFlags: []string{"tls", "accesspoint=" + accessPointID},
						},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
				TargetPath: targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", []string{"accesspoint=" + accessPointID, "tls"}},
			mountSuccess:  true,
		},
		{
			name: "fail: conflicting access point in volume handle and mount options",
			req: &csi.NodePublishVolumeRequest{
				VolumeId: volumeId + "::" + accessPointID,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{
							MountFlags: []string{"tls", "accesspoint=fsap-deadbeef"},
						},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
				TargetPath: targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Found conflicting access point IDs in mountOptions (fsap-deadbeef) and volumeHandle (fsap-abcd1234)",
			},
		},
		{
			name: "fail: too many fields in volume handle",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + ":/a/b/::four!",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "volume ID 'fs-abc123:/a/b/::four!' is invalid: Expected at most three fields separated by ':'",
			},
		},
		{
			name: "fail: missing target path",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Target path not provided",
			},
		},
		{
			name: "fail: missing volume capability",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:   volumeId,
				TargetPath: targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Volume capability not provided",
			},
		},
		{
			name: "fail: unsupported volume capability",
			req: &csi.NodePublishVolumeRequest{
				VolumeId: volumeId,
				VolumeCapability: &csi.VolumeCapability{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
				TargetPath: targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Volume capability not supported",
			},
		},
		{
			name: "fail: mounter failed to MakeDir",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{}, // Signal MakeDir failure
			expectError: errtyp{
				code:    "Internal",
				message: `Could not create dir "/target/path": failed to MakeDir`,
			},
		},
		{
			name: "fail: mounter failed to Mount",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: true,
			mountArgs:     []interface{}{volumeId + ":/", targetPath, "efs", gomock.Any()},
			mountSuccess:  false,
			expectError: errtyp{
				code:    "Internal",
				message: `Could not mount "fs-abc123:/" at "/target/path": failed to Mount`,
			},
		},
		{
			name: "fail: unsupported volume context",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId,
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
				VolumeContext:    map[string]string{"asdf": "qwer"},
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Volume context property asdf not supported",
			},
		},
		{
			name: "fail: invalid filesystem ID",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         "invalid-id",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "volume ID 'invalid-id' is invalid: Expected a file system ID of the form 'fs-...'",
			},
		},
		{
			name: "fail: invalid access point ID",
			req: &csi.NodePublishVolumeRequest{
				VolumeId:         volumeId + "::invalid-id",
				VolumeCapability: stdVolCap,
				TargetPath:       targetPath,
			},
			expectMakeDir: false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "volume ID 'fs-abc123::invalid-id' has an invalid access point ID 'invalid-id': Expected it to be of the form 'fsap-...'",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockMounter, driver, ctx := setup(mockCtrl)

			if tc.expectMakeDir {
				var err error
				// If not expecting mount, it's because mkdir errored
				if len(tc.mountArgs) == 0 {
					err = fmt.Errorf("failed to MakeDir")
				}
				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(err)
			}
			if len(tc.mountArgs) != 0 {
				var err error
				if !tc.mountSuccess {
					err = fmt.Errorf("failed to Mount")
				}
				mockMounter.EXPECT().Mount(tc.mountArgs[0], tc.mountArgs[1], tc.mountArgs[2], tc.mountArgs[3]).Return(err)
			}

			ret, err := driver.NodePublishVolume(ctx, tc.req)
			testResult(t, "NodePublishVolume", ret, err, tc.expectError)
		})
	}
}

func TestNodeUnpublishVolume(t *testing.T) {
	testCases := []struct {
		name                string
		req                 *csi.NodeUnpublishVolumeRequest
		expectGetDeviceName bool
		getDeviceNameReturn []interface{}
		expectUnmount       bool
		unmountReturn       error
		expectError         errtyp
	}{
		{
			name: "success: normal",
			req: &csi.NodeUnpublishVolumeRequest{
				VolumeId:   volumeId,
				TargetPath: targetPath,
			},
			expectGetDeviceName: true,
			getDeviceNameReturn: []interface{}{"", 1, nil},
			expectUnmount:       true,
			unmountReturn:       nil,
		},
		{
			name: "success: unpublish with already unmounted target",
			req: &csi.NodeUnpublishVolumeRequest{
				VolumeId:   volumeId,
				TargetPath: targetPath,
			},
			expectGetDeviceName: true,
			getDeviceNameReturn: []interface{}{"", 0, nil},
			// NUV returns early if the refcount is zero
			expectUnmount: false,
		},
		{
			name: "fail: targetPath is missing",
			req: &csi.NodeUnpublishVolumeRequest{
				VolumeId: volumeId,
			},
			expectGetDeviceName: false,
			expectUnmount:       false,
			expectError: errtyp{
				code:    "InvalidArgument",
				message: "Target path not provided",
			},
		},
		{
			name: "fail: mounter failed to umount",
			req: &csi.NodeUnpublishVolumeRequest{
				VolumeId:   volumeId,
				TargetPath: targetPath,
			},
			expectGetDeviceName: true,
			getDeviceNameReturn: []interface{}{"", 1, nil},
			expectUnmount:       true,
			unmountReturn:       fmt.Errorf("Unmount failed"),
			expectError: errtyp{
				code:    "Internal",
				message: `Could not unmount "/target/path": Unmount failed`,
			},
		},
		{
			name: "fail: mounter failed to GetDeviceName",
			req: &csi.NodeUnpublishVolumeRequest{
				VolumeId:   volumeId,
				TargetPath: targetPath,
			},
			expectGetDeviceName: true,
			getDeviceNameReturn: []interface{}{"", 1, fmt.Errorf("GetDeviceName failed")},
			expectUnmount:       false,
			expectError: errtyp{
				code:    "Internal",
				message: "failed to check if volume is mounted: GetDeviceName failed",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockMounter, driver, ctx := setup(mockCtrl)

			if tc.expectGetDeviceName {
				mockMounter.EXPECT().
					GetDeviceName(targetPath).
					Return(tc.getDeviceNameReturn[0], tc.getDeviceNameReturn[1], tc.getDeviceNameReturn[2])
			}
			if tc.expectUnmount {
				mockMounter.EXPECT().Unmount(targetPath).Return(tc.unmountReturn)
			}

			ret, err := driver.NodeUnpublishVolume(ctx, tc.req)
			testResult(t, "NodeUnpublishVolume", ret, err, tc.expectError)
		})
	}
}
