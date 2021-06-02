package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/junglemc/JungleTree/internal/player"
	"github.com/junglemc/JungleTree/net"
	. "github.com/junglemc/JungleTree/net/codec"
	. "github.com/junglemc/JungleTree/packet"
	"github.com/junglemc/JungleTree/pkg"
	"log"
)

func playPluginMessage(c *net.Client, p net.Packet) (err error) {
	pkt := p.(ServerboundPluginMessagePacket)

	if pkt.Channel.Prefix() == "minecraft" && pkt.Channel.Name() == "brand" {
		buf := bytes.NewBuffer(pkt.Data)
		brand := ReadString(buf)

		if onlinePlayer, ok := player.GetOnlinePlayer(c); ok {
			onlinePlayer.ClientBrand = brand

			if pkg.Config().DebugMode {
				log.Printf("Client brand for %s: %s\n", c.Profile.Name, onlinePlayer.ClientBrand)
			}
		}
	}
	return
}

func playClientSettings(c *net.Client, p net.Packet) (err error) {
	pkt := p.(ServerboundClientSettingsPacket)

	onlinePlayer, ok := player.GetOnlinePlayer(c)
	if !ok {
		return
	}
	onlinePlayer.Locale = pkt.Locale
	onlinePlayer.ViewDistance = pkt.ViewDistance
	onlinePlayer.ChatMode = pkt.ChatMode
	onlinePlayer.ChatColorsEnabled = pkt.ChatColorsEnabled
	onlinePlayer.MainHand = pkt.MainHand

	if pkg.Config().DebugMode {
		data, _ := json.MarshalIndent(onlinePlayer, "", "  ")
		log.Printf("Client settings for %s:\n%s\n", c.Profile.Name, string(data))
	}

	itemChange := &ClientboundHeldItemChangePacket{
		Slot: byte(onlinePlayer.Hotbar.SelectedIndex),
	}
	return c.Send(itemChange)
}
