package config

// DeviceDimensions holds the dimensions and DPR for a specific device.
type DeviceDimensions struct {
	PortraitW  int
	PortraitH  int
	LandscapeW int
	LandscapeH int
	DPR        float64
}

// deviceNames lists all supported devices in order. The index of the device
// corresponds to the DeviceType integer used in params.RenderParam.
var deviceNames = []string{
	"Oneline", "iPhone 16", "iPhone 16 Plus", "iPhone 16 Pro", "iPhone 16 Pro Max",
	"Pixel 7", "Pixel 8", "Pixel 8 Pro", "Pixel 8a",
	"Pixel 9", "Pixel 9 Pro", "Pixel 9 Pro XL", "Pixel 10",
	"Samsung Galaxy A55", "Samsung Galaxy S8+", "Samsung Galaxy S20 Ultra", "Samsung Galaxy A51/71",
	"Facebook on Android", "Surface Duo",
	"Pixel 9 Pro Fold", "Galaxy Z Fold 5", "Galaxy Z Fold 6",
	"Asus Zenbook Fold", "iPad Mini", "iPad Pro 13",
	"Surface Pro 10", "Surface Pro 7", "Nest Hub Max",
}

// DeviceTypeMap maps the integer DeviceType (index) to the device dimensions.
var DeviceTypeMap = map[int]DeviceDimensions{
	0:  {},
	1:  {393, 852, 852, 393, 3.0},
	2:  {430, 932, 932, 430, 3.0},
	3:  {402, 874, 874, 402, 3.0},
	4:  {440, 956, 956, 440, 3.0},
	5:  {412, 892, 892, 412, 2.625},
	6:  {412, 892, 892, 412, 2.625},
	7:  {448, 968, 968, 448, 3.25},
	8:  {393, 852, 852, 393, 3.0},
	9:  {412, 915, 915, 412, 2.625},
	10: {412, 915, 915, 412, 3.125},
	11: {448, 968, 968, 448, 3.25},
	12: {412, 915, 915, 412, 2.625},
	13: {384, 832, 832, 384, 2.0},
	14: {360, 740, 740, 360, 4.0},
	15: {412, 915, 915, 412, 3.5},
	16: {412, 915, 915, 412, 2.5},
	17: {360, 640, 640, 360, 3.0},
	18: {540, 720, 720, 540, 2.5},
	19: {992, 862, 862, 992, 3.0},
	20: {768, 952, 952, 768, 2.6},
	21: {968, 848, 848, 968, 2.6},
	22: {1280, 1710, 1710, 1280, 1.5},
	23: {744, 1133, 1133, 744, 2.0},
	24: {1032, 1376, 1376, 1032, 2.0},
	25: {912, 1368, 1368, 912, 2.0},
	26: {912, 1368, 1368, 912, 2.0},
	27: {1280, 800, 1280, 800, 2.0},
}
