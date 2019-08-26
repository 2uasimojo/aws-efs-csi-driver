# v0.2.0

## Changelog

### Notable changes
* Combine manifest files ([#35](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/35), [@leakingtapan](https://github.com/leakingtapan))
* Add example stateful sets ([#43](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/43), [@leakingtapan](https://github.com/leakingtapan))
* Added flag for version information output ([#44](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/44), [@djcass44](https://github.com/djcass44))
* fix namespace in csi-node clusterrole ([#47](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/47), [@d-nishi](https://github.com/d-nishi))
* Update to CSI v1.1.0 ([#48](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/48), [@wongma7](https://github.com/wongma7))
* Add support for 'path' field in volumeContext ([#52](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/52), [@wongma7](https://github.com/wongma7))
* Replace deprecated Recycle policy with Retain ([#53](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/53), [@wongma7](https://github.com/wongma7))
* Add sanity test ([#54](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/54), [@wongma7](https://github.com/wongma7))
* Run upstream e2e tests  ([#55](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/55), [@wongma7](https://github.com/wongma7))
* Add linux nodeSelector to manifest.yaml ([#61](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/61), [@wongma7](https://github.com/wongma7))
* Add liveness probe ([#62](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/62), [@wongma7](https://github.com/wongma7))
* Add example for volume path ([#65](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/65), [@leakingtapan](https://github.com/leakingtapan))
* Upgrade to golang 1.12 ([#70](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/70), [@wongma7](https://github.com/wongma7))

# v0.1.0
[Documentation](https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/v0.1.0/docs/README.md)

filename  | sha512 hash
--------- | ------------
[v0.1.0.zip](https://github.com/kubernetes-sigs/aws-efs-csi-driver/archive/v0.1.0.zip) | `b2ac6ccedfbd40f7a47ed1c14fb9bc16742592f03c3f51e26ef5f72ed2f97718cae32dca998304f5773c3b0d3df100942817d55bbb09cbd2226a51000cfc1505`
[v0.1.0.tar.gz](https://github.com/kubernetes-sigs/aws-efs-csi-driver/archive/v0.1.0.tar.gz) | `1db081d96906ae07a868cbcf3e3902fe49c44f219966c1f5ba5a8beabd9311e42cae57ff1884edf63b936afce128b113ed94d85afc2e2955dedb81ece99f72dc`

## Changelog

### Notable changes
* Multiple README updates and example updates
* Switch to use klog for logging ([#20](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/20), [@leakingtapan](https://github.com/leakingtapan/))
* Update README and add more examples ([#18](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/18), [@leakingtapan](https://github.com/leakingtapan/))
* Update manifest files ([#12](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/12), [@leakingtapan](https://github.com/leakingtapan/))
* Add sample manifest for multiple pod RWX scenario ([#9](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/9), [@leakingtapan](https://github.com/leakingtapan/))
* Update travis with code verification ([#8](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/8), [@leakingtapan](https://github.com/leakingtapan/))
* Implement mount options support ([#5](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/5), [@leakingtapan](https://github.com/leakingtapan/))
* Update logging format of the driver ([#4](https://github.com/kubernetes-sigs/aws-efs-csi-driver/pull/4), [@leakingtapan](https://github.com/leakingtapan/))
* Implement node service for EFS driver  ([bca5d36](https://github.com/kubernetes-sigs/aws-efs-csi-driver/commit/bca5d36), [@leakingtapan](https://github.com/leakingtapan/))
