#!/usr/bin/env python3

# fmt: off

import sys
sys.dont_write_bytecode = True

from scripts import utils as x
# ----------------------------

import datetime as dt
import os

from pathlib import Path

if __name__ == "__main__":
    _build_entry = (Path(x.PROJ_ROOT) / 'build.py').resolve().absolute().as_posix()
    help_str  = f'Usage: [python3] {_build_entry} --help\n'
    x.print_stderr(help_str[:-1])
    sys.exit(1)



BUILD_CMD  = 'go'
BUILD_ENV  = os.environ.copy()
GO_LDFLAGS = '-v'

_target_platform = ''
_target_archlibc = ''
_target_optimization = 'default'

_pkg_inst_dir = ''

_extra_args_build: list[str] = []

def module_init(env: dict) -> list:
    global BUILD_CMD; \
        BUILD_CMD = env['GOCMD_EXEC']
    global GO_LDFLAGS; \
        GO_LDFLAGS = env['GO_LDFLAGS']
    global _target_optimization; \
        _target_optimization = env['OPTIMIZATION']
    global _pkg_inst_dir; \
        _pkg_inst_dir = env['PKG_INST_DIR']
    global _extra_args_build; \
        _extra_args_build = env['EXTRA_ARGS']

    BUILD_ENV.update({
        **env['GOENV_BULD']
    })

    return [
        _build_step_00,
    ]



def _build_step_00():
    x._util_func__subprocess(env=BUILD_ENV, args=[BUILD_CMD, 'env'])

    _go_module = x._util_func__subprocess(collect_stdout=True, cwd=x.PROJ_ROOT, env=BUILD_ENV, args=[BUILD_CMD, 'list', '-m']).strip()
    _go_suffix = x._util_func__subprocess(collect_stdout=True, cwd=x.PROJ_ROOT, env=BUILD_ENV, args=[BUILD_CMD, 'env', 'GOEXE']).strip()
    x._util_put_pkg_version_desc(x._util_func__subprocess(cwd=x.PROJ_ROOT, collect_stdout=True, args=['git', 'describe', '--always', '--abbrev=7']))

    _go_ldflags  = GO_LDFLAGS
    _go_ldflags += f" -X '{_go_module}/system.version={x._util_get_pkg_version_desc()}'"

    _build_datetime = dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%d %H:%M:%S %Z")
    _go_ldflags += f" -X '{_go_module}/system.datetime={_build_datetime}'"

    _go_ldflags += f" -X '{_go_module}/system.flavor={_target_optimization}'"


    args = [BUILD_CMD, 'build', #'-x',
        '-o', f"{_pkg_inst_dir}/{_go_module}{_go_suffix}",
        '-ldflags', _go_ldflags,
    ]
    x._util_func__subprocess(cwd=x.PROJ_ROOT, env=BUILD_ENV, args=[*args, *_extra_args_build, '.'])
