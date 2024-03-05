package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

func settingsTab(pomoConfig *PomoConfig) *fyne.Container {
	timingConfigSelection := widget.NewSelect(maps.Keys(pomoConfig.TimingConfigs), nil)
	timingConfigSelection.SetSelectedIndex(0)

	newTimingConfigEntry := widget.NewEntry()
	newTimingConfigEntry.SetPlaceHolder("Config name")

	newTimingConfigSaveButton := widget.NewButton("Save", nil)
	newTimingConfigContainer := container.NewBorder(nil, nil, nil, newTimingConfigSaveButton, newTimingConfigEntry)

	newTimingConfigSaveButton.OnTapped = func() {
		pomoConfig.NewTimingConfig(newTimingConfigEntry.Text)
		timingConfigSelection.SetOptions(maps.Keys(pomoConfig.TimingConfigs))
		timingConfigSelection.Refresh()
		timingConfigSelection.SetSelected(newTimingConfigEntry.Text)
		newTimingConfigContainer.Hide()
	}

	newTimingConfigContainer.Hide()

	timingConfigSelectionNewButton := widget.NewButton("New", nil)
	timingConfigSelectionNewButton.OnTapped = func() {
		newTimingConfigContainer.Show()
	}
	timingConfigSelectionDeleteButton := widget.NewButton("Delete", nil)
	timingConfigSelectionDeleteButton.OnTapped = func() {
		prevSelection := timingConfigSelection.Selected
		if prevSelection != DefaultTimingConfigName {
			pomoConfig.DeleteTimingConfig(prevSelection)
			timingConfigSelection.SetOptions(maps.Keys(pomoConfig.TimingConfigs))
			timingConfigSelection.SetSelected(DefaultTimingConfigName)
		}
	}
	timingConfigButtonsContainer := container.NewHBox(timingConfigSelectionNewButton, timingConfigSelectionDeleteButton)
	timingConfigSelectionContainer := container.NewBorder(nil, nil, nil, timingConfigButtonsContainer, timingConfigSelection)

	settingsForm := createSettingsForm(pomoConfig, DefaultTimingConfigName)

	timingConfigSettingsContainer := container.NewVBox(
		timingConfigSelectionContainer,
		newTimingConfigContainer,
		settingsForm,
	)

	// TODO: See if Form widget's optional buttons could be used
	saveButton := widget.NewButton("Save", nil)
	saveButton.OnTapped = func() {
		pomoConfig.Save()
	}

	timingConfigSelection.OnChanged = func(s string) {
		settingsForm = createSettingsForm(pomoConfig, s)
		// TODO: Do this without specifying the index
		timingConfigSettingsContainer.Objects[2] = settingsForm
	}

	return container.NewVBox(
		timingConfigSettingsContainer,
		// layout.NewSpacer(),
		saveButton,
	)
}

func createSettingsForm(pomoConfig *PomoConfig, timingConfigName string) *widget.Form {
	log.Debugf("createSettingsForm: %s", timingConfigName)
	return widget.NewForm(
		widget.NewFormItem("Work Duration", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigName].WorkDuration))),
		widget.NewFormItem("Break Duration", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigName].BreakDuration))),
		widget.NewFormItem("Long Break Duration", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigName].LongBreakDuration))),
		widget.NewFormItem("Work Iterations", widget.NewEntryWithData(binding.IntToString(pomoConfig.TimingConfigs[timingConfigName].WorkIterations))),
		widget.NewFormItem("Auto Start Next", widget.NewCheckWithData("", pomoConfig.TimingConfigs[timingConfigName].AutoStartNext)),

		widget.NewFormItem("Work Finished Sound", widget.NewEntryWithData(pomoConfig.TimingConfigs[timingConfigName].WorkFinishedSound)),
		widget.NewFormItem("Break Finished Sound", widget.NewEntryWithData(pomoConfig.TimingConfigs[timingConfigName].BreakFinishedSound)),
		widget.NewFormItem("Play Work Finished Sound", widget.NewCheckWithData("", pomoConfig.TimingConfigs[timingConfigName].PlayWorkFinishedSound)),
		widget.NewFormItem("Play Break Finished Sound", widget.NewCheckWithData("", pomoConfig.TimingConfigs[timingConfigName].PlayBreakFinishedSound)),
	)
}
