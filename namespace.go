package uuid

var (
	// DNSNamespace a fully qualified domain name
	DNSNamespace UUID

	// URLNamespace is a URL
	URLNamespace UUID

	// IODNamespace is an ISO OID
	IODNamespace UUID

	// X500Namespace is an X.500 DN
	X500Namespace UUID
)

// Namespaces taken from Appendix C
// https://tools.ietf.org/html/rfc4122#appendix-C
func initNamespace() error {
	var err error

	if DNSNamespace, err = FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8"); err != nil {
		return err
	}

	if URLNamespace, err = FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8"); err != nil {
		return err
	}

	if IODNamespace, err = FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8"); err != nil {
		return err
	}

	if X500Namespace, err = FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8"); err != nil {
		return err
	}

	return nil
}
