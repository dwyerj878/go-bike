# Bike-Go

## Overview
The application is intended to take an activity file and generate statistics such as normalized power and training effects based purely on the athlete and sample data.

Initially it is using json format as exported by Golden Cheetah but the intent is to support multiple inout formats and create a consistent output that will work in garmin connect, golden cheetah etc.

## Technology
* The initial implementation is using go 1.23.4

## TODO
- [x] load json gc file
- [ ] open zip file
- [ ] load garmin .fit file
- [x] athlete data
- [x] normalized power
- [ ] training effect
- [ ] time in zones
- [ ] graph output (https://github.com/gonum/plot) ?
- [ ] implement test coverage (https://github.com/marketplace/actions/go-test-coverage)
- [ ] implement trends and multi activity analysis

## NP Formula
- Step 1: Calculate the rolling average power with a window size of 30 seconds. Start at 30s and calculate the average of the previous 30s and repeat this for every second.
- Step 2: Take each value from step one and take this value to the fourth power (multiply this number by itself four times).
- Step 3: Calculate the average of values from the previous step.
- Step 4: Take the fourth root of the average from the previous step — this value gives us the normalized power.
