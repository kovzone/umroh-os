// catalog_version_messages.go — hand-written proto message types for
// package versioning RPC (BL-CAT-013).
//
// Run `make genpb` to regenerate from catalog.proto once protoc is available,
// then delete this file.
package pb

// GetPackageVersionRequest carries the package_id to version.
type GetPackageVersionRequest struct {
	PackageId string
}

func (x *GetPackageVersionRequest) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}

// GetPackageVersionResponse carries the computed version hash and metadata.
type GetPackageVersionResponse struct {
	PackageId string
	// Version is the MD5 hex of (package_id + "|" + name + "|" + status + "|" + updated_at_unix).
	Version   string
	UpdatedAt string
}

func (x *GetPackageVersionResponse) GetPackageId() string {
	if x == nil {
		return ""
	}
	return x.PackageId
}
func (x *GetPackageVersionResponse) GetVersion() string {
	if x == nil {
		return ""
	}
	return x.Version
}
func (x *GetPackageVersionResponse) GetUpdatedAt() string {
	if x == nil {
		return ""
	}
	return x.UpdatedAt
}
