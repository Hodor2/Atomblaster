package audio

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// SoundType constants for different sound effects
const (
	ShootSound = iota
	HitSound
	PickupSound
	DoorSound
	DeathSound
	DashSound
)

// AudioSystem manages all game sounds
type AudioSystem struct {
	ShootSound      rl.Sound
	HitSound        rl.Sound
	PickupSound     rl.Sound
	DoorSound       rl.Sound
	DeathSound      rl.Sound
	DashSound       rl.Sound
	BackgroundMusic rl.Music
}

// NewAudioSystem initializes the audio system and loads all sounds
func NewAudioSystem() *AudioSystem {
	rl.InitAudioDevice()
	
	system := &AudioSystem{
		ShootSound:  GenerateSound(800, 0.1, 0.8),   // 800 Hz for shooting
		HitSound:    GenerateSound(400, 0.15, 0.7),  // 400 Hz for hit events
		PickupSound: GenerateSound(1000, 0.2, 0.9),  // 1000 Hz for gun pickup
		DoorSound:   GenerateSound(600, 0.3, 0.8),   // 600 Hz for door unlock
		DeathSound:  GenerateSound(200, 0.5, 0.9),   // 200 Hz for death
		DashSound:   GenerateSound(1200, 0.1, 0.7),  // 1200 Hz for dash
	}
	
	// Set the volume for all sounds
	rl.SetSoundVolume(system.ShootSound, 0.7)
	rl.SetSoundVolume(system.HitSound, 0.7)
	rl.SetSoundVolume(system.PickupSound, 0.8)
	rl.SetSoundVolume(system.DoorSound, 0.8)
	rl.SetSoundVolume(system.DeathSound, 0.8)
	rl.SetSoundVolume(system.DashSound, 0.7)
	
	// Try to load background music if available
	// system.BackgroundMusic = rl.LoadMusicStream("assets/background_music.mp3")
	// rl.PlayMusicStream(system.BackgroundMusic)
	// rl.SetMusicVolume(system.BackgroundMusic, 0.5)
	
	return system
}

// PlaySound plays a sound effect based on the sound type
func (as *AudioSystem) PlaySound(soundType int) {
	switch soundType {
	case ShootSound:
		rl.PlaySound(as.ShootSound)
	case HitSound:
		rl.PlaySound(as.HitSound)
	case PickupSound:
		rl.PlaySound(as.PickupSound)
	case DoorSound:
		rl.PlaySound(as.DoorSound)
	case DeathSound:
		rl.PlaySound(as.DeathSound)
	case DashSound:
		rl.PlaySound(as.DashSound)
	}
}

// Update updates the audio system (for streaming music)
func (as *AudioSystem) Update() {
	if as.BackgroundMusic.CtxData != nil {
		rl.UpdateMusicStream(as.BackgroundMusic)
	}
}

// Cleanup releases audio resources
func (as *AudioSystem) Cleanup() {
	// Unload all sound resources
	rl.UnloadSound(as.ShootSound)
	rl.UnloadSound(as.HitSound)
	rl.UnloadSound(as.PickupSound)
	rl.UnloadSound(as.DoorSound)
	rl.UnloadSound(as.DeathSound)
	rl.UnloadSound(as.DashSound)
	
	// Unload music if loaded
	if as.BackgroundMusic.CtxData != nil {
		rl.UnloadMusicStream(as.BackgroundMusic)
	}
	
	rl.CloseAudioDevice()
}