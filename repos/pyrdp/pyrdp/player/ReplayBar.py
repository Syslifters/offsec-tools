#
# This file is part of the PyRDP project.
# Copyright (C) 2018-2024 GoSecure Inc.
# Licensed under the GPLv3 or later.
#
import datetime
import logging

from PySide6.QtCore import Qt, Signal
from PySide6.QtWidgets import QCheckBox, QHBoxLayout, QLabel, QSizePolicy, QSlider, QSpacerItem, \
    QVBoxLayout, QWidget

from pyrdp.logging import LOGGER_NAMES
from pyrdp.player.SeekBar import SeekBar
from pyrdp.ui import PlayPauseButton


class ReplayBar(QWidget):
    """
    Widget that contains the play/pause button, the seek bar and the speed slider.
    """
    play = Signal()
    pause = Signal()
    seek = Signal(float)
    speedChanged = Signal(int)

    def __init__(self, duration: float, referenceTime: int, parent: QWidget = None):
        super().__init__(parent)

        self.log = logging.getLogger(LOGGER_NAMES.PLAYER)
        
        # absolute capture starting point
        self.referenceTime = referenceTime
        
        # pretty display capture duration
        self.duration = duration
        self.duration_hours = int(self.duration // 3600)
        remaining_seconds = int(self.duration % 3600)
        minutes = int(remaining_seconds // 60)
        seconds = int(remaining_seconds % 60)
        if self.duration_hours > 0:
            self.duration_pretty = "{:d}:{:02d}:{:02d}".format(
                                    self.duration_hours, minutes, seconds)
        else:
            self.duration_pretty = "{:02d}:{:02d}".format(minutes, seconds)

        self.button = PlayPauseButton()
        self.button.setMaximumWidth(100)
        self.button.clicked.connect(self.onButtonClicked)

        self.scaleCheckbox = QCheckBox("Scale to window")

        self.timeSlider = SeekBar()
        self.timeSlider.setMinimum(0)
        self.timeSlider.setMaximum(int(self.duration * 1000))
        self.timeSlider.valueChanged.connect(self.onSeek)

        self.timeLabel = QLabel(self.formatTimeLabel(0))
        self.absoluteDateTimeLabel = QLabel(self.formatAbsoluteDateTimeLabel(0))

        self.speedLabel = QLabel("Speed: 1x")

        self.speedSlider = QSlider(Qt.Horizontal)
        self.speedSlider.setMaximumWidth(300)
        self.speedSlider.setMinimum(1)
        self.speedSlider.setMaximum(10)
        self.speedSlider.valueChanged.connect(self.onSpeedChanged)

        vertical = QVBoxLayout()

        horizontal = QHBoxLayout()
        horizontal.addWidget(self.speedLabel)
        horizontal.addWidget(self.speedSlider)
        horizontal.addWidget(self.scaleCheckbox)
        horizontal.addItem(QSpacerItem(20, 40, QSizePolicy.Expanding, QSizePolicy.Expanding))
        vertical.addLayout(horizontal)

        horizontal = QHBoxLayout()
        horizontal.addWidget(self.button)
        horizontal.addWidget(self.timeSlider)
        horizontal.addWidget(self.timeLabel)
        horizontal.addWidget(self.absoluteDateTimeLabel)
        vertical.addLayout(horizontal)

        self.setLayout(vertical)
        self.setGeometry(0, 0, 80, 60)

    def onButtonClicked(self):
        if self.button.playing:
            self.log.debug("Play clicked")
            self.play.emit()
        else:
            self.log.debug("Pause clicked")
            self.pause.emit()

    def onSeek(self):
        time = self.timeSlider.value() / 1000.0
        self.log.debug("Seek to %(arg1)d seconds", {"arg1": time})
        self.seek.emit(time)

    def onSpeedChanged(self):
        speed = self.speedSlider.value()
        self.log.debug("Slider changed value: %(arg1)d", {"arg1": speed})
        self.speedLabel.setText("Speed: {}x".format(speed))
        self.speedChanged.emit(speed)

    def onTimeChanged(self, currentTime: float):
        if currentTime >= self.duration:
            currentTime = self.duration
        
        self.timeLabel.setText(self.formatTimeLabel(currentTime))
        self.absoluteDateTimeLabel.setText(self.formatAbsoluteDateTimeLabel(currentTime))    

    def formatTimeLabel(self, current: float) -> str:
        hours = int(current // 3600)
        remaining_seconds = int(current % 3600)
        minutes = int(remaining_seconds // 60)
        seconds = int(remaining_seconds % 60)
        if self.duration_hours > 0:
            return "{:d}:{:02d}:{:02d} / {}".format(
                hours, minutes, seconds, self.duration_pretty)
        else:
            return "{:02d}:{:02d} / {}".format(minutes, seconds, self.duration_pretty)
        
    def formatAbsoluteDateTimeLabel(self, current: int) -> str:
        current = current * 1000
        epoch_time = (self.referenceTime + current)
        formatted_time = datetime.datetime.fromtimestamp(epoch_time / 1000, tz=datetime.timezone.utc).strftime('%Y-%m-%d %H:%M:%S UTC')

        return formatted_time
