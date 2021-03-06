package uuid

import (
	"testing"
)

const (
	testSize = 100000
)

func devNull(i interface{}) {}

func TestInsertTimestamp(t *testing.T) {

	b := make([]byte, 8)
	full := uint64(0xFFFFFFFFFFFFFFFF)
	insertTimestamp(b, uint64(full))

	for _, v := range b {
		if v != 0xFF {
			t.Error("Not all bytes were set in Insert Timestamp", v)
		}
	}
}

func TestVersion(t *testing.T) {
	uuid := UUID{}

	for i := 0; i < 6; i++ {
		uuid.version(uint8(i))
		if uuid[6] != uint8(i<<4) {
			t.Error("Version is not correct:", uuid[6], "should be:", uint8(i<<4))
		}
	}
}

func TestVariant(t *testing.T) {
	uuid := UUID{}

	for i := uint8(0); i < 0xFF; i++ {
		uuid[8] = i
		uuid.variant(rfc4122)

		if uuid[8] < 0x0F || uuid[8] > 0xBF {
			t.Error("Varient is not correct", uuid[8], "at", i)
		}
	}

}

func TestFromStringBadFormat(t *testing.T) {

	t.Parallel()

	tests := []struct {
		uuid string
	}{
		{
			uuid: "6ba7b814-9dad-61d1-80b4-00c04fd430c8", // wrong version
		},
		{
			uuid: "6ba7b814-9dad-11d1-30b4-00c04fd430c8", // wrong variant
		},
	}

	for _, test := range tests {
		_, err := FromString(test.uuid)
		if err == nil {
			t.Error("FromString did not detect bad uuid String")
		}
	}

}

func TestFromBytesBadFormat(t *testing.T) {
	b := make([]byte, 16)
	_, err := FromBytes(b)

	if err != ErrUUIDFormat {
		t.Error("FromBytes did not detect bad uuid String")
	}
}

func TestFromBytesWrongLen(t *testing.T) {

	uuid := NewV1()
	b := make([]byte, 10)
	copy(b, uuid[:])

	_, err := FromBytes(b)

	if err != ErrUUIDSize {
		t.Error("FromBytes did not detect wrong length")
	}
}

func TestRegexV1(t *testing.T) {

	for i := 0; i < testSize; i++ {
		uuid := NewV1()

		if !uuidRegex.MatchString(uuid.String()) {
			t.Error("v1 does not pass regex test", uuid.String())
		}
	}
}

func TestMutexV1(t *testing.T) {

	for i := 0; i < testSize/10; i++ {
		go func() {
			NewV1()
		}()
	}
}

func TestCollisionV1(t *testing.T) {
	uuids := make(map[UUID]uint8)

	for i := 0; i < testSize; i++ {
		uuid := NewV1()

		_, ok := uuids[uuid]

		if ok == true { //collision
			t.Error("Collision V1:", uuid.String())
		} else {
			uuids[uuid] = 0
		}
	}
}

func TestRegexV2(t *testing.T) {

	for i := 0; i < testSize; i++ {
		uuid := NewV2()
		if !uuidRegex.MatchString(uuid.String()) {
			t.Error("V2 does not pass regex test", uuid.String())
		}
	}
}

func TestMutexV2(t *testing.T) {

	for i := 0; i < testSize/10; i++ {
		go func() {
			NewV2()
		}()
	}
}

func TestCollisionV2(t *testing.T) {
	uuids := make(map[UUID]uint8)

	for i := 0; i < testSize; i++ {
		uuid := NewV2()

		_, ok := uuids[uuid]

		if ok == true { //collision
			t.Error("Collision V2:", uuid.String())
		} else {
			uuids[uuid] = 0
		}
	}
}

func TestRegexV3(t *testing.T) {

	uuid, err := NewV3(DNSNamespace, "google")

	if err != nil {
		t.Error("V3 error", err)
	}

	if !uuidRegex.MatchString(uuid.String()) {
		t.Error("V3 does not pass regex test", uuid.String())
	}
}

// V3 requires collisions for same name and namespace
// See https://tools.ietf.org/html/rfc4122#section-4.3
func TestCollisionV3(t *testing.T) {

	uuid, u1err := NewV3(URLNamespace, "google")
	uuid2, u2err := NewV3(URLNamespace, "google")

	if u1err != nil || u2err != nil {
		t.Error("V3 error", u1err, u2err)
	}

	if uuid.String() != uuid.String() {
		t.Error("V3 does not pass collision", uuid.String(), uuid2.String())
	}

}

func TestRegexV4(t *testing.T) {

	for i := 0; i < testSize; i++ {
		uuid := NewV4()

		if !uuidRegex.MatchString(uuid.String()) {
			t.Error("V4 does not pass regex test", uuid.String())
		}
	}
}

func TestMutexV4(t *testing.T) {

	for i := 0; i < testSize/10; i++ {
		go func() {
			NewV4()
		}()
	}
}

func TestCollisionV4(t *testing.T) {
	uuids := make(map[UUID]uint8)

	for i := 0; i < testSize; i++ {
		uuid := NewV4()

		_, ok := uuids[uuid]

		if ok == true { //collision
			t.Error("Collision V4:", uuid.String())
		} else {
			uuids[uuid] = 0
		}
	}
}

func TestRegexV5(t *testing.T) {

	uuid, err := NewV5(DNSNamespace, "google")

	if err != nil {
		t.Error("V5 error", err)
	}

	if !uuidRegex.MatchString(uuid.String()) {
		t.Error("V5 does not pass regex test", uuid.String())
	}
}

// V5 requires collisions for same name and namespace
// See https://tools.ietf.org/html/rfc4122#section-4.3
func TestCollisionV5(t *testing.T) {

	uuid, u1err := NewV5(DNSNamespace, "google")
	uuid2, u2err := NewV5(DNSNamespace, "google")

	if u1err != nil || u2err != nil {
		t.Error("V3 error", u1err, u2err)
	}

	if uuid.String() != uuid.String() {
		t.Error("V5 does not pass collision", uuid.String(), uuid2.String())
	}

}

func TestClockSeqInit(t *testing.T) {
	var cs uint16
	var dup int

	for i := 0; i < testSize; i++ {
		temp := clockSeqInit()

		if cs == temp {
			dup++
		}

		cs = temp
	}
	if dup > 10 {
		t.Error("Clock Sequence is not random")
	}
}

func BenchmarkV1(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uuid := NewV1()
		devNull(uuid)
	}
}

func BenchmarkV2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uuid := NewV2()
		devNull(uuid)
	}
}

func BenchmarkV3(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uuid, _ := NewV3(DNSNamespace, "name")
		devNull(uuid)
	}
}

func BenchmarkV4(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uuid := NewV4()
		devNull(uuid)
	}
}

func BenchmarkV5(b *testing.B) {
	for n := 0; n < b.N; n++ {
		uuid, _ := NewV5(DNSNamespace, "name")
		devNull(uuid)
	}
}
