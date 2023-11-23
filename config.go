package corewar

var IND_SIZE = 2
var REG_SIZE = 4
var DIR_SIZE = REG_SIZE

var REG_CODE = 1
var DIR_CODE = 2
var IND_CODE = 3

var MAX_PLAYERS = 4

const MEM_SIZE = 4 * 1024

var IDX_MOD = MEM_SIZE / 8
var PLAYER_MAX_SIZE = MEM_SIZE / 6

var COMMENT_CHAR = '#'
var LABEL_CHAR = ':'
var DIRECT_CHAR = '%'
var SEPARATOR_CHAR = ','

var LABEL_CHARS = "abcdefghijklmnopqrstuvwxyz_0123456789"

var NAME_CMD_STRING = ".name"
var DESCRIPTION_CMD_STRING = ".description"

const REG_NUMBER = 16

var CYCLE_TO_DIE = 1536
var CYCLE_DELTA = 50
var NBR_LIVE = 21
var MAX_CHECKS = 10

var PROG_NAME_LENGTH = 128
var DESCRIPTION_LENGTH = 2048
var COREWAR_EXEC_SIGNATURE = 0xea83f3
