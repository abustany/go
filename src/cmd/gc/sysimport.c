char*	sysimport =
	"package sys\n"
	"type sys._esys_002 {}\n"
	"type sys.any 24\n"
	"type sys._esys_003 *sys.any\n"
	"type sys._osys_281 {_esys_279 sys._esys_003}\n"
	"type sys.uint32 6\n"
	"type sys._isys_283 {_esys_280 sys.uint32}\n"
	"type sys._esys_001 (sys._esys_002 sys._osys_281 sys._isys_283)\n"
	"var !sys.mal sys._esys_001\n"
	"type sys._esys_005 {}\n"
	"type sys._esys_006 {}\n"
	"type sys._esys_007 {}\n"
	"type sys._esys_004 (sys._esys_005 sys._esys_006 sys._esys_007)\n"
	"var !sys.breakpoint sys._esys_004\n"
	"type sys._esys_009 {}\n"
	"type sys._esys_010 {}\n"
	"type sys.int32 5\n"
	"type sys._isys_289 {_esys_288 sys.int32}\n"
	"type sys._esys_008 (sys._esys_009 sys._esys_010 sys._isys_289)\n"
	"var !sys.panicl sys._esys_008\n"
	"type sys._esys_012 {}\n"
	"type sys._esys_013 {}\n"
	"type sys.bool 12\n"
	"type sys._isys_294 {_esys_293 sys.bool}\n"
	"type sys._esys_011 (sys._esys_012 sys._esys_013 sys._isys_294)\n"
	"var !sys.printbool sys._esys_011\n"
	"type sys._esys_015 {}\n"
	"type sys._esys_016 {}\n"
	"type sys.float64 10\n"
	"type sys._isys_299 {_esys_298 sys.float64}\n"
	"type sys._esys_014 (sys._esys_015 sys._esys_016 sys._isys_299)\n"
	"var !sys.printfloat sys._esys_014\n"
	"type sys._esys_018 {}\n"
	"type sys._esys_019 {}\n"
	"type sys.int64 7\n"
	"type sys._isys_304 {_esys_303 sys.int64}\n"
	"type sys._esys_017 (sys._esys_018 sys._esys_019 sys._isys_304)\n"
	"var !sys.printint sys._esys_017\n"
	"type sys._esys_021 {}\n"
	"type sys._esys_022 {}\n"
	"type sys._esys_023 25\n"
	"type sys.string *sys._esys_023\n"
	"type sys._isys_309 {_esys_308 sys.string}\n"
	"type sys._esys_020 (sys._esys_021 sys._esys_022 sys._isys_309)\n"
	"var !sys.printstring sys._esys_020\n"
	"type sys._esys_025 {}\n"
	"type sys._esys_026 {}\n"
	"type sys.uint8 2\n"
	"type sys._esys_027 *sys.uint8\n"
	"type sys._isys_314 {_esys_313 sys._esys_027}\n"
	"type sys._esys_024 (sys._esys_025 sys._esys_026 sys._isys_314)\n"
	"var !sys.printpointer sys._esys_024\n"
	"type sys._esys_029 {}\n"
	"type sys._osys_321 {_esys_318 sys.string}\n"
	"type sys._isys_323 {_esys_319 sys.string _esys_320 sys.string}\n"
	"type sys._esys_028 (sys._esys_029 sys._osys_321 sys._isys_323)\n"
	"var !sys.catstring sys._esys_028\n"
	"type sys._esys_031 {}\n"
	"type sys._osys_331 {_esys_328 sys.int32}\n"
	"type sys._isys_333 {_esys_329 sys.string _esys_330 sys.string}\n"
	"type sys._esys_030 (sys._esys_031 sys._osys_331 sys._isys_333)\n"
	"var !sys.cmpstring sys._esys_030\n"
	"type sys._esys_033 {}\n"
	"type sys._osys_342 {_esys_338 sys.string}\n"
	"type sys._isys_344 {_esys_339 sys.string _esys_340 sys.int32 _esys_341 sys.int32}\n"
	"type sys._esys_032 (sys._esys_033 sys._osys_342 sys._isys_344)\n"
	"var !sys.slicestring sys._esys_032\n"
	"type sys._esys_035 {}\n"
	"type sys._osys_353 {_esys_350 sys.uint8}\n"
	"type sys._isys_355 {_esys_351 sys.string _esys_352 sys.int32}\n"
	"type sys._esys_034 (sys._esys_035 sys._osys_353 sys._isys_355)\n"
	"var !sys.indexstring sys._esys_034\n"
	"type sys._esys_037 {}\n"
	"type sys._osys_362 {_esys_360 sys.string}\n"
	"type sys._isys_364 {_esys_361 sys.int64}\n"
	"type sys._esys_036 (sys._esys_037 sys._osys_362 sys._isys_364)\n"
	"var !sys.intstring sys._esys_036\n"
	"type sys._esys_039 {}\n"
	"type sys._osys_371 {_esys_368 sys.string}\n"
	"type sys._esys_040 *sys.uint8\n"
	"type sys._isys_373 {_esys_369 sys._esys_040 _esys_370 sys.int32}\n"
	"type sys._esys_038 (sys._esys_039 sys._osys_371 sys._isys_373)\n"
	"var !sys.byteastring sys._esys_038\n"
	"type sys._esys_042 {}\n"
	"type sys._esys_043 <>\n"
	"type sys._osys_382 {_esys_378 sys._esys_043}\n"
	"type sys._esys_044 *sys.uint8\n"
	"type sys._esys_045 *sys.uint8\n"
	"type sys._ssys_389 {}\n"
	"type sys._esys_046 *sys._ssys_389\n"
	"type sys._isys_384 {_esys_379 sys._esys_044 _esys_380 sys._esys_045 _esys_381 sys._esys_046}\n"
	"type sys._esys_041 (sys._esys_042 sys._osys_382 sys._isys_384)\n"
	"var !sys.mkiface sys._esys_041\n"
	"type sys._esys_048 {}\n"
	"type sys._osys_393 {_esys_392 sys.int32}\n"
	"type sys._esys_049 {}\n"
	"type sys._esys_047 (sys._esys_048 sys._osys_393 sys._esys_049)\n"
	"var !sys.argc sys._esys_047\n"
	"type sys._esys_051 {}\n"
	"type sys._osys_397 {_esys_396 sys.int32}\n"
	"type sys._esys_052 {}\n"
	"type sys._esys_050 (sys._esys_051 sys._osys_397 sys._esys_052)\n"
	"var !sys.envc sys._esys_050\n"
	"type sys._esys_054 {}\n"
	"type sys._osys_402 {_esys_400 sys.string}\n"
	"type sys._isys_404 {_esys_401 sys.int32}\n"
	"type sys._esys_053 (sys._esys_054 sys._osys_402 sys._isys_404)\n"
	"var !sys.argv sys._esys_053\n"
	"type sys._esys_056 {}\n"
	"type sys._osys_410 {_esys_408 sys.string}\n"
	"type sys._isys_412 {_esys_409 sys.int32}\n"
	"type sys._esys_055 (sys._esys_056 sys._osys_410 sys._isys_412)\n"
	"var !sys.envv sys._esys_055\n"
	"type sys._esys_058 {}\n"
	"type sys._osys_419 {_esys_416 sys.float64 _esys_417 sys.int32}\n"
	"type sys._isys_421 {_esys_418 sys.float64}\n"
	"type sys._esys_057 (sys._esys_058 sys._osys_419 sys._isys_421)\n"
	"var !sys.frexp sys._esys_057\n"
	"type sys._esys_060 {}\n"
	"type sys._osys_428 {_esys_425 sys.float64}\n"
	"type sys._isys_430 {_esys_426 sys.float64 _esys_427 sys.int32}\n"
	"type sys._esys_059 (sys._esys_060 sys._osys_428 sys._isys_430)\n"
	"var !sys.ldexp sys._esys_059\n"
	"type sys._esys_062 {}\n"
	"type sys._osys_438 {_esys_435 sys.float64 _esys_436 sys.float64}\n"
	"type sys._isys_440 {_esys_437 sys.float64}\n"
	"type sys._esys_061 (sys._esys_062 sys._osys_438 sys._isys_440)\n"
	"var !sys.modf sys._esys_061\n"
	"type sys._esys_064 {}\n"
	"type sys._osys_447 {_esys_444 sys.bool}\n"
	"type sys._isys_449 {_esys_445 sys.float64 _esys_446 sys.int32}\n"
	"type sys._esys_063 (sys._esys_064 sys._osys_447 sys._isys_449)\n"
	"var !sys.isInf sys._esys_063\n"
	"type sys._esys_066 {}\n"
	"type sys._osys_456 {_esys_454 sys.bool}\n"
	"type sys._isys_458 {_esys_455 sys.float64}\n"
	"type sys._esys_065 (sys._esys_066 sys._osys_456 sys._isys_458)\n"
	"var !sys.isNaN sys._esys_065\n"
	"type sys._esys_068 {}\n"
	"type sys._osys_464 {_esys_462 sys.float64}\n"
	"type sys._isys_466 {_esys_463 sys.int32}\n"
	"type sys._esys_067 (sys._esys_068 sys._osys_464 sys._isys_466)\n"
	"var !sys.Inf sys._esys_067\n"
	"type sys._esys_070 {}\n"
	"type sys._osys_471 {_esys_470 sys.float64}\n"
	"type sys._esys_071 {}\n"
	"type sys._esys_069 (sys._esys_070 sys._osys_471 sys._esys_071)\n"
	"var !sys.NaN sys._esys_069\n"
	"type sys._esys_073 {}\n"
	"type sys._esys_075 [sys.any] sys.any\n"
	"type sys._esys_074 *sys._esys_075\n"
	"type sys._osys_474 {hmap sys._esys_074}\n"
	"type sys._isys_476 {keysize sys.uint32 valsize sys.uint32 keyalg sys.uint32 valalg sys.uint32 hint sys.uint32}\n"
	"type sys._esys_072 (sys._esys_073 sys._osys_474 sys._isys_476)\n"
	"var !sys.newmap sys._esys_072\n"
	"type sys._esys_077 {}\n"
	"type sys._osys_485 {val sys.any}\n"
	"type sys._esys_079 [sys.any] sys.any\n"
	"type sys._esys_078 *sys._esys_079\n"
	"type sys._isys_487 {hmap sys._esys_078 key sys.any}\n"
	"type sys._esys_076 (sys._esys_077 sys._osys_485 sys._isys_487)\n"
	"var !sys.mapaccess1 sys._esys_076\n"
	"type sys._esys_081 {}\n"
	"type sys._osys_493 {val sys.any pres sys.bool}\n"
	"type sys._esys_083 [sys.any] sys.any\n"
	"type sys._esys_082 *sys._esys_083\n"
	"type sys._isys_495 {hmap sys._esys_082 key sys.any}\n"
	"type sys._esys_080 (sys._esys_081 sys._osys_493 sys._isys_495)\n"
	"var !sys.mapaccess2 sys._esys_080\n"
	"type sys._esys_085 {}\n"
	"type sys._esys_086 {}\n"
	"type sys._esys_088 [sys.any] sys.any\n"
	"type sys._esys_087 *sys._esys_088\n"
	"type sys._isys_502 {hmap sys._esys_087 key sys.any val sys.any}\n"
	"type sys._esys_084 (sys._esys_085 sys._esys_086 sys._isys_502)\n"
	"var !sys.mapassign1 sys._esys_084\n"
	"type sys._esys_090 {}\n"
	"type sys._esys_091 {}\n"
	"type sys._esys_093 [sys.any] sys.any\n"
	"type sys._esys_092 *sys._esys_093\n"
	"type sys._isys_508 {hmap sys._esys_092 key sys.any val sys.any pres sys.bool}\n"
	"type sys._esys_089 (sys._esys_090 sys._esys_091 sys._isys_508)\n"
	"var !sys.mapassign2 sys._esys_089\n"
	"type sys._esys_095 {}\n"
	"type sys._esys_097 1 sys.any\n"
	"type sys._esys_096 *sys._esys_097\n"
	"type sys._osys_515 {hchan sys._esys_096}\n"
	"type sys._isys_517 {elemsize sys.uint32 elemalg sys.uint32 hint sys.uint32}\n"
	"type sys._esys_094 (sys._esys_095 sys._osys_515 sys._isys_517)\n"
	"var !sys.newchan sys._esys_094\n"
	"type sys._esys_099 {}\n"
	"type sys._esys_100 {}\n"
	"type sys._esys_101 {}\n"
	"type sys._esys_098 (sys._esys_099 sys._esys_100 sys._esys_101)\n"
	"var !sys.gosched sys._esys_098\n"
	"type sys._esys_103 {}\n"
	"type sys._esys_104 {}\n"
	"type sys._esys_105 {}\n"
	"type sys._esys_102 (sys._esys_103 sys._esys_104 sys._esys_105)\n"
	"var !sys.goexit sys._esys_102\n"
	"type sys._esys_107 {}\n"
	"type sys._osys_529 {_esys_526 sys.string _esys_527 sys.bool}\n"
	"type sys._isys_531 {_esys_528 sys.string}\n"
	"type sys._esys_106 (sys._esys_107 sys._osys_529 sys._isys_531)\n"
	"var !sys.readfile sys._esys_106\n"
	"type sys._esys_109 {}\n"
	"type sys._osys_540 {_esys_535 sys.int32 _esys_536 sys.int32}\n"
	"type sys._esys_110 *sys.uint8\n"
	"type sys._isys_542 {_esys_537 sys._esys_110 _esys_538 sys.int32 _esys_539 sys.int32}\n"
	"type sys._esys_108 (sys._esys_109 sys._osys_540 sys._isys_542)\n"
	"var !sys.bytestorune sys._esys_108\n"
	"type sys._esys_112 {}\n"
	"type sys._osys_553 {_esys_548 sys.int32 _esys_549 sys.int32}\n"
	"type sys._isys_555 {_esys_550 sys.string _esys_551 sys.int32 _esys_552 sys.int32}\n"
	"type sys._esys_111 (sys._esys_112 sys._osys_553 sys._isys_555)\n"
	"var !sys.stringtorune sys._esys_111\n"
	"type sys._esys_114 {}\n"
	"type sys._esys_115 {}\n"
	"type sys._isys_562 {_esys_561 sys.int32}\n"
	"type sys._esys_113 (sys._esys_114 sys._esys_115 sys._isys_562)\n"
	"var !sys.exit sys._esys_113\n"
	"))\n"
;
