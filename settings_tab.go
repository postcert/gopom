package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/exp/maps"
)

func settingsTab(pomoConfig *PomoConfig) *fyne.Container {
	timingConfigSelection := widget.NewSelect(maps.Keys(pomoConfig.TimingConfigs), nil)
	timingConfigSelection.SetSelectedIndex(0)

	timingConfigSelectionNewButton := widget.NewButton("New", nil)
	timingConfigSelectionContainer := container.NewHBox(timingConfigSelection, timingConfigSelectionNewButton)

	settingsForm := widget.NewForm(
		widget.NewFormItem("Work Duration", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigSelection.Selected].WorkDuration))),
		widget.NewFormItem("Break Duration", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigSelection.Selected].BreakDuration))),
		widget.NewFormItem("Work Iterations", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigSelection.Selected].WorkIterations))),
		widget.NewFormItem("Auto Start Next", widget.NewCheckWithData("", pomoConfig.TimingConfigs[timingConfigSelection.Selected].AutoStartNext)),
	)
	autoStartNext := widget.NewCheckWithData("Auto Start Next", pomoConfig.TimingConfigs[timingConfigSelection.Selected].AutoStartNext)

	timingConfigSettingsContainer := container.NewVBox(
		timingConfigSelectionContainer,
		settingsForm,
	)

	// TODO: See if Form widget's optional buttons could be used
	saveButton := widget.NewButton("Save", nil)
	saveButton.OnTapped = func() {
		pomoConfig.Save()
	}

	return container.NewVBox(
		timingConfigSettingsContainer,
		settingsForm,
		saveButton,
		autoStartNext,
	)
}
