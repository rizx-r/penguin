package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {
	if size <= 0 {
		size = 250
	}
	return &Bitmap{bits: make([]byte, size), size: size << 3}
}

func (b *Bitmap) Set(id string) {
	// id在哪个bit
	idx := hash(id) % b.size
	// 计算在哪个byte
	byteIdx := idx >> 3
	// 在这个byte中的那个bit位置
	bitIdx := idx & 7 /*idx % 8*/
	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	idx := hash(id) % b.size
	byteIdx := idx >> 3
	bitIdx := idx & 7
	return (b.bits[byteIdx] & (1 << bitIdx)) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

func LoadBitmap(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{bits: bits, size: len(bits) << 3}
}

func hash(id string) int {
	// BKDR哈希算法
	seed := 131313 // 31 131 1313 13131 131313, etc
	hs := 0
	for _, c := range id {
		hs = hs*seed + int(c)
	}
	return hs & 0x7fffffff
}
