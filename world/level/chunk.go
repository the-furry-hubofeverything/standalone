package level

import (
	"github.com/junglemc/JungleTree/world/blocks"
)

const (
	chunkSize   = 16
	biomeBlocks = 4
)

type Chunk struct {
	World     *World
	X         byte
	Z         byte
	sections  []ChunkSection
	heightMap [chunkSize * chunkSize]int32 // signed integer, future proofing for the new chunk format in 1.17.?
	biomes    [chunkSize * chunkSize * biomeBlocks]int32
}

// TODO: Call before sending the chunk across the network, of course.
func (c *Chunk) Update() {
	c.updateHeightMap()
	c.updateBiomes()

	if c.sections == nil {
		return
	}
	for _, section := range c.sections {
		section.updatePalette()
	}
}

func (c *Chunk) BlockAt(x, y, z int32) *blocks.Block {
	modX := x % chunkSectionSize
	modY := y % chunkSectionSize
	modZ := z % chunkSectionSize

	return c.sections[(y-modY)/chunkSectionSize].BlockAt(modX, modY, modZ)
}

func (c *Chunk) SetBlock(x, y, z int32, block *blocks.Block) {
	modX := x % chunkSectionSize
	modY := y % chunkSectionSize
	modZ := z % chunkSectionSize

	c.sections[(y-modY)/chunkSectionSize].SetBlock(modX, modY, modZ, block)
}

func (c *Chunk) HeightMap() (heightMap [chunkSize * chunkSize]int32) {
	return c.heightMap
}

func (c *Chunk) updateHeightMap() {
	if c.sections == nil {
		return
	}

	pos := 0

	for i := len(c.sections) - 1; i >= 0; i-- {
		for z := 0; z < chunkSectionSize; z++ {
			for x := 0; x < chunkSectionSize; x++ {
				y, ok := c.sections[i].HighestBlockAt(x, z)
				if !ok {
					continue
				}
				c.heightMap[pos] = y
			}
		}
	}
}

// TODO: Test function fills whole array with int32(127)
func (c *Chunk) updateBiomes() {
	// TODO: Just fills with void type for now - need to calculate it from a smoothed voroni diagram
	for x := 0; x < chunkSize; x++ {
		for z := 0; z < chunkSize; z++ {
			for y := 0; y < int(c.World.Height); y++ {
				i := ((y>>2)&63)<<4 | ((z>>2)&3)<<2 | ((x >> 2) & 3)
				c.biomes[i] = 127 // TODO: Biome ID for void
			}
		}
	}
}
