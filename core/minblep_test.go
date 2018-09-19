//-----------------------------------------------------------------------------
/*

MinBLEP testing

*/
//-----------------------------------------------------------------------------

package core

import (
	"math"
	"testing"
)

//-----------------------------------------------------------------------------
// Floating Point Comparisons
// See: http://floating-point-gui.de/errors/NearlyEqualsTest.java

const minNormal = 2.2250738585072014E-308 // 2**-1022

func equal(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)
	if a == 0 || b == 0 || diff < minNormal {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * minNormal)
	}
	// use relative error
	return diff/math.Min((absA+absB), math.MaxFloat64) < epsilon
}

func convertToFloat64(ibuf []uint64) (fbuf []float64) {
	fbuf = make([]float64, len(ibuf))
	for i := 0; i < len(ibuf); i++ {
		fbuf[i] = math.Float64frombits(ibuf[i])
	}
	return
}

//-----------------------------------------------------------------------------

func TestSinc(t *testing.T) {
	xi := []uint64{
		0x401b370b70000000, 0xc000e61350000000, 0x4016a5df50000000, 0x4017e00d50000000, 0x40207744e8000000, 0xc0183225e0000000, 0xc00a5d46b0000000, 0x40157557a0000000,
		0xc011c72c78000000, 0x3ff1453880000000, 0xbfdcee8880000000, 0x40049e8d60000000, 0xc005a26d80000000, 0x3fd1273600000000, 0x402216d520000000, 0x4020a5d678000000,
		0x4005b6c0a0000000, 0x4011623de0000000, 0xc01cabfac4000000, 0x40011d7200000000, 0xc023591507400000, 0xc01491add4000000, 0xc01d057f54000000, 0x4018558a30000000,
		0xc01b773650000000, 0xbfffb2a2e0000000, 0xc01d9de44c000000, 0xc01f4b9884000000, 0x4023f4fca0000000, 0xc0168a193c000000, 0x3fd08db200000000, 0x401b2104b0000000,
	}
	yi := []uint64{
		0x3f9bb33863696c8e, 0x3faaab09f2644feb, 0xbfa922fd328f6486, 0xbf75602fd8d4a8f0, 0x3f9a74ad4f70e59b, 0x3f8083ebb9dd5c54, 0xbfb3cc7e3feab1f3, 0xbfabac0f43fa8d45,
		0x3fb20e74aa03a301, 0xbfb2a2ed8410b62b, 0x3fe646f69cb12cde, 0x3fbeaf7495c04ec6, 0x3fb822b32aa56d4a, 0x3fec59c21dbf46f5, 0xbf74212736956efd, 0x3fa0a896d9fb4eff,
		0x3fb779ae6f4ea939, 0x3fb098bbef6c001d, 0xbf96e5200580b508, 0x3fb02673a37406e6, 0xbf9cc8bc0b51a306, 0xbf9b655d4b764b3d, 0xbfa0262af7e3bd62, 0x3f8bccd0629d8684,
		0x3f935777bb2dd1ce, 0xbf83838da923116a, 0xbfa5058613e743d4, 0xbf95e622bf5fd743, 0xbf61a548c4ea4598, 0xbfaa5d832ab46d5f, 0x3fec97d03b59fb74, 0x3f9e5f4a79c9fcaa,
	}
	x := convertToFloat64(xi)
	y := convertToFloat64(yi)
	for i := 0; i < len(xi); i++ {
		result := Sinc(x[i])
		if !equal(y[i], result, 1e-12) {
			t.Logf("for x %f expected %.20f, actual %.20f\n", x[i], y[i], result)
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestBlackman(t *testing.T) {
	wi := []uint64{
		0xbc70000000000000, 0x3f6ebbca113e7b00, 0x3f90038676002b96, 0x3fa3267257557c80, 0x3fb24b8125b1878f, 0x3fbecb17f77fc969, 0x3fc7c345979582a8, 0x3fd1262c9b07126b,
		0x3fd76833ef642891, 0x3fde6c9a193c12d9, 0x3fe2eb583c2a7379, 0x3fe6978ef74c6408, 0x3fe9f7c0c3383945, 0x3fecc8a026e5e20a, 0x3feeceb09ecbef93, 0x3fefdd91ca280c48,
		0x3fefdd91ca280c4c, 0x3feeceb09ecbefa2, 0x3fecc8a026e5e221, 0x3fe9f7c0c3383962, 0x3fe6978ef74c6429, 0x3fe2eb583c2a739a, 0x3fde6c9a193c1318, 0x3fd76833ef6428d1,
		0x3fd1262c9b07129a, 0x3fc7c345979582fc, 0x3fbecb17f77fc9ee, 0x3fb24b8125b187f9, 0x3fa3267257557cf6, 0x3f90038676002c32, 0x3f6ebbca113e7d40, 0xbc70000000000000,
	}
	w := convertToFloat64(wi)
	n := len(w)
	y := BlackmanWindow(n)
	for i := 0; i < n; i++ {
		if !equal(y[i], w[i], 1e-13) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, w[i], y[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestDFT(t *testing.T) {

	realTimei := []uint64{
		0x400205ba60000000, 0xc0105145b8000000, 0x4006022400000000, 0x3fdf166b00000000, 0xbfc06d7300000000, 0x4022e93768000000, 0xc010994188000000, 0x4015b567f0000000,
		0x3fe11de500000000, 0x401597d5f0000000, 0xbfffed4580000000, 0x401f5285e0000000, 0xc01155b6d8000000, 0xc0079b4e40000000, 0x40189e32d0000000, 0x4020c2d4d0000000,
		0xc02135b4b9000000, 0x4021f91be8000000, 0x3fe0a31500000000, 0xc0208ec9c4000000, 0xc0189f75dc000000, 0x400a1dc620000000, 0x401f37f6d0000000, 0xc0082d5880000000,
		0xc0216ee2bb000000, 0xc02332f6c8000000, 0xbfeb122580000000, 0xc02179e60c000000, 0xc014f006cc000000, 0x4022d34b10000000, 0x4020169c60000000, 0x401c12d650000000,
	}
	imagTimei := []uint64{
		0xc012aaaf70000000, 0x3fe9725600000000, 0xc003f785b0000000, 0x4014d1e4d0000000, 0x3fd00b9900000000, 0x400ad5f740000000, 0x3fe43a6600000000, 0xc0226dc4ef800000,
		0xbff3f4ba60000000, 0x402145fdb0000000, 0x40213b7e08000000, 0x4011ad1a90000000, 0xc01141abc0000000, 0x4013152ee0000000, 0x4006658860000000, 0xc0075a2aa0000000,
		0x400e0ecd60000000, 0xc01ab8d958000000, 0xbff32aa360000000, 0x401e67f0c0000000, 0x401a5609e0000000, 0xc00b256410000000, 0xc015aebb58000000, 0x401f784460000000,
		0xc007f13f30000000, 0x400dddffa0000000, 0x4022423c10000000, 0x3ffc5d66c0000000, 0x40092b2c80000000, 0x401cb1b0e0000000, 0xbff35740a0000000, 0x4020f57368000000,
	}
	realFreqi := []uint64{
		0x4030e36bc8000000, 0x3ff1c43d82f0ed34, 0x4005ebbb8d1af263, 0x4024e51710b69832, 0x403eddd4d2e90096, 0xc0486a090e270dc7, 0xc0431b413c53f7d7, 0x4033f3e7cd1e83fa,
		0xc04abece6f1fffaf, 0x403925007fb3b762, 0xc03811a6de63c6ce, 0xbfc709e4c5a1f3c0, 0xc02df92c7c87e75d, 0x4042a3f6405ff870, 0xc012d0d9dd9b9e15, 0x4028235452588dec,
		0xc0448c0315fffefc, 0x40377e841518c04e, 0x3ff23680336b7c8c, 0x4013a079220dd538, 0x401913c6d45bfffb, 0x404ff87dc8225f4f, 0xc04283da1e0272c5, 0x4042ff4b24f544cf,
		0xc045ec7e76e00092, 0xc037d9871d906980, 0x404dc28e76c4d8b2, 0xc0401eb998bb2687, 0xc030e02649bc0369, 0x40368857ffde5892, 0x404ba0886809dc50, 0x4033f2602ded7844,
	}
	imagFreqi := []uint64{
		0x404a7c9216200000, 0xc040a958e3b80605, 0xc0218a1f657e26e4, 0x403ead0fea47a803, 0x403820608dacea44, 0x3ffec2a055fc7518, 0x4001c6ae5cd31480, 0x4048fa70042be868,
		0xc0351f126e000286, 0xc04dcca2d7858db1, 0x404037de84804d9c, 0x4041a4dd9a5997c8, 0x4031da6fc6491984, 0xc02e45cac13d64e8, 0xbffbbe5e60c3c740, 0xc04efe3b83bb960a,
		0xc03e3b11b44001e4, 0x4021e141fe57f797, 0xc018ac7a12a47f3e, 0xc050d10720333039, 0xc025d243b359d00e, 0xc037892873ecc5f2, 0x4043493df5e5ac1f, 0xc0205bc24271ef46,
		0x3fdca5bd7ffefdd0, 0xc00906141f82299d, 0xc0314fe5919864da, 0xc0322bafd06071c2, 0xc052955cd6924706, 0xc0318b33a194d34e, 0xc029341462b2f1ac, 0x404726a89f4558d2,
	}

	realTime := convertToFloat64(realTimei)
	imagTime := convertToFloat64(imagTimei)
	realFreq := convertToFloat64(realFreqi)
	imagFreq := convertToFloat64(imagFreqi)

	x, y := DFT(realTime, imagTime)
	n := len(realTime)
	for i := 0; i < n; i++ {
		if !equal(x[i], realFreq[i], 1e-11) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, realFreq[i], x[i])
			t.Error("FAIL")
		}
		if !equal(y[i], imagFreq[i], 1e-10) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, imagFreq[i], y[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestInverseDFT(t *testing.T) {

	realFreqi := []uint64{
		0xc0004008c0000000, 0x40192e6d00000000, 0x400d799740000000, 0x4020705a98000000, 0xbfd6697600000000, 0xc016bbe7a0000000, 0x4022029588000000, 0x4020ce1cf0000000,
		0xc01c2fec40000000, 0x401e7c2720000000, 0x400692a980000000, 0xbff5c66000000000, 0x400322ac00000000, 0xc01183e738000000, 0x4016e152b0000000, 0xc00ece86a0000000,
		0xbff0f30240000000, 0xc015e95614000000, 0xc018ff5268000000, 0xc011e6b6b0000000, 0x3ff20fdf40000000, 0xbffab83700000000, 0xc01a6e7264000000, 0x402045ac08000000,
		0xc01fbf0ddc000000, 0xc01de9fa1c000000, 0xbfb7539000000000, 0x4014d68850000000, 0x402363db58000000, 0x40216670e0000000, 0x400d82de40000000, 0xc002b09b10000000,
	}
	imagFreqi := []uint64{
		0x4013fb4ec0000000, 0xc005038a80000000, 0xc010779880000000, 0xc0156b48a0000000, 0x3ffb094c80000000, 0xc014726d5c000000, 0xc01bcf0e98000000, 0x40129266e0000000,
		0xc01df64628000000, 0x40177a4600000000, 0xc01adf3148000000, 0x40139b1000000000, 0xc02104d09c000000, 0x40220110c8000000, 0xc021e619b2000000, 0x3fdb99e100000000,
		0xc019e734ac000000, 0xc014cb85cc000000, 0x4017d2e760000000, 0x40129cc310000000, 0x40090cd8c0000000, 0x4022b23a88000000, 0x40065036e0000000, 0x4014c75ec0000000,
		0xc02042c28b000000, 0xc01d3532d0000000, 0x3fd9de7100000000, 0xc020dee724000000, 0xc021342892000000, 0xc017a0a9f8000000, 0xbfe8b0de80000000, 0x401992fdb0000000,
	}
	realTimei := []uint64{
		0x3fe6a2b852000000, 0x3ffcd5be951f392f, 0xbff26bece78ab738, 0x3fe1be846c5647e7, 0xbfd2228450761506, 0xbfe78c19ebff215f, 0x3fd0dd02d7d5c6a8, 0x3fecab3655b2327f,
		0xbfb69cec10000b1c, 0xbfe9b9ec07842c44, 0x3fe636a15182c49a, 0x3fe2492828a767c9, 0x3f7a13ad8e865f00, 0xbfedc417f96f43ec, 0x3ffd26877cb806a6, 0xbffb0848da29d7e4,
		0xbfd237c85c0005ee, 0x3fe2b381a5684e1e, 0xbf7973940847cb50, 0xbfe556ec3e48e72c, 0xc001e85ddbf13d95, 0x3ff834f541788bef, 0x3ff2b35a756924e2, 0xbfb1967bb14c4514,
		0xbfef88888e00018d, 0xbff266dc73a28e60, 0xbfe9545d0a5cc856, 0xbfecad057ac8bd9b, 0xbff5a1ccf78e8528, 0xbfe5d617fa601e77, 0x3fef2aa1b7d2ccb0, 0x3ff376331a6f7a43,
	}
	imagTimei := []uint64{
		0x3ff29ff798400000, 0xbfa9945d755917cc, 0x3fdc02df5dfaf857, 0xbff8385fe50f9adb, 0x3fd910c8c72da744, 0xbffed2756091bb48, 0xbff1f030003bb24f, 0xbff7c18fcafe33ec,
		0x3fef0c93488000bc, 0x3fe23a2bbc59de55, 0xbff47e87ed8479b3, 0xbfdda55c550b7516, 0xbfc29efab8ceb4ed, 0x3ff37d092e3b7cf5, 0xbff50cf7e5792b76, 0xbfedc91da7409a18,
		0x3ffd2010d63fffee, 0xbff7354cccd2505c, 0xbfc7bd57e304cf70, 0x3fea359f7d642790, 0x3fe45ed82569272e, 0x3fe1cb0a2efa594c, 0x3fc83704063142c3, 0xbff065d5e2f67d36,
		0xbfd1f7803f000389, 0x3fd30b798953d8e0, 0x3fe596bd12ccaa84, 0x3fde755b10827ffa, 0xbfd201ff1998b237, 0xc0026fe13a95de51, 0xbfeddf02c8228925, 0x3ff7a2232914c425,
	}

	realTime := convertToFloat64(realTimei)
	imagTime := convertToFloat64(imagTimei)
	realFreq := convertToFloat64(realFreqi)
	imagFreq := convertToFloat64(imagFreqi)
	x, y := InverseDFT(realFreq, imagFreq)
	n := len(realFreq)
	for i := 0; i < n; i++ {
		if !equal(x[i], realTime[i], 1e-11) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, realTime[i], x[i])
			t.Error("FAIL")
		}
		if !equal(y[i], imagTime[i], 1e-4) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, imagTime[i], y[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestRealCepstrum(t *testing.T) {

	signali := []uint64{
		0x3ff7764280000000, 0x4014724bb0000000, 0xc021ec2582000000, 0xc01b601c2c000000, 0x4023ffef20000000, 0xc017a7599c000000, 0x401f324aa0000000, 0xc01df667d4000000,
		0x4023e97628000000, 0xc021d6734f000000, 0x401da4a7f0000000, 0xc0211b5a65000000, 0xc023d56297900000, 0x4020ec3a50000000, 0x3ffe0ba540000000, 0xc01991f9d8000000,
		0xc01af31118000000, 0xc001545e50000000, 0x40208564b0000000, 0x4019935b70000000, 0xc0068b7430000000, 0x3ff0cb9340000000, 0x3ff96ae780000000, 0xbfee59f9c0000000,
		0x400dfb6360000000, 0xc02003af8d000000, 0x3fe3b79400000000, 0x40149560a0000000, 0xc00f501210000000, 0x4023b06b58000000, 0x3ff8a17940000000, 0x401e3587c0000000,
	}
	realCepstrumi := []uint64{
		0x400ae7ebe4000000, 0xbfcbee3860268bb8, 0x3f67dda2ddca15c0, 0xbfb8ff728a42059d, 0xbfa1c18b4db2921e, 0x3f982f455432530c, 0xbf7bcfe1d5995ba0, 0xbfa3d572f9536913,
		0xbfa0c17a0000005f, 0x3fc8ccd5ad06a379, 0x3f905b96ff2a59ce, 0xbf95ff54caf2ea94, 0xbfb1d6b15926b8a4, 0x3fa220b5bcdbad64, 0xbfa73f1372beaaff, 0xbfc63b90c9c9167d,
		0xbfba2a4e80000000, 0xbfc63b90c9c915a1, 0xbfa73f1372bea37e, 0x3fa220b5bcdbb258, 0xbfb1d6b15926b571, 0xbf95ff54caf2da9c, 0x3f905b96ff2a6964, 0x3fc8ccd5ad06a4ef,
		0xbfa0c179fffffcc8, 0xbfa3d572f953643a, 0xbf7bcfe1d5991628, 0x3f982f4554326328, 0xbfa1c18b4db28488, 0xbfb8ff728a41fc30, 0x3f67dda2ddcbc500, 0xbfcbee3860267cc5,
	}

	signal := convertToFloat64(signali)
	realCepstrum := convertToFloat64(realCepstrumi)

	x := RealCepstrum(signal)
	n := len(signal)
	for i := 0; i < n; i++ {
		if !equal(x[i], realCepstrum[i], 1e-6) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, realCepstrum[i], x[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------

func TestMinBLEP(t *testing.T) {

	minBLEPi := []uint64{
		0xbf63f31459fe7913, 0xbf504befd8694319, 0x3f85bcf96c5402ca, 0x3fa63c7c3e149ccf, 0x3fbd080f3459c989, 0x3fce5436686d1d8b, 0x3fdb0199f282b97b, 0x3fe511a3fa904c7b,
		0x3fed417151b27a8e, 0x3ff2459267092b7a, 0x3ff4b818d92e0115, 0x3ff58c06c6562e56, 0x3ff4dc055e2c8adf, 0x3ff342af9a0631f8, 0x3ff1979dcf8aad22, 0x3ff090e1643eb41e,
		0x3ff0797fe63516dc, 0x3ff121881af7ccd5, 0x3ff20b435f9236ce, 0x3ff2b752818e5225, 0x3ff2e3b26747e233, 0x3ff29e9a83e2be12, 0x3ff22b704a6543ad, 0x3ff1d1faa2bf37a3,
		0x3ff1b7464ef4e049, 0x3ff1d2f2a151afe0, 0x3ff1fec5626f21e0, 0x3ff21053cca06ae2, 0x3ff1ec6bda2ac09a, 0x3ff18db5984fecb7, 0x3ff101e4d36d92a2, 0x3ff063bd419491ad,
		0x3fefa95831b85d68, 0x3feee974538aa429, 0x3feeb007089508e4, 0x3feefa4838a783cd, 0x3fef963d4e5d4d0e, 0x3ff01ad626d765da, 0x3ff048f8d2fc048c, 0x3ff046a989c43980,
		0x3ff01ee554dbd383, 0x3fefdaf54a81fee1, 0x3fef9b5189eb8eaa, 0x3fef96f313324d71, 0x3fefc06f2d1de742, 0x3feff4c884bac2d9, 0x3ff00a445baad4b1, 0x3ff00a6d1db446db,
		0x3ff0000000000000,
	}

	minBLEP := convertToFloat64(minBLEPi)

	x := GenerateMinBLEP(6, 4)
	n := len(minBLEP)
	for i := 0; i < n; i++ {
		if !equal(x[i], minBLEP[i], 1e-6) {
			t.Logf("for i %d expected %.20f, actual %.20f\n", i, minBLEP[i], x[i])
			t.Error("FAIL")
		}
	}
}

//-----------------------------------------------------------------------------
