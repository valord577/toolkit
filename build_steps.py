#!/usr/bin/env python3

# fmt: off

import datetime as dt
import os
import shlex
import shutil


_env: dict = {}
_ctx: dict = {
    'PKG_INST_STRIP': '',
    'BUILD_ENV': os.environ.copy(),

    'EXTRA_ARGS_CONFIGURE': [],
}

def module_init(env: dict) -> list:
    global _env; _env = env
    return [
        _build_step_00,
    ]



def _build_step_00():
    _build_env = {
        **_ctx['BUILD_ENV'],
        **_env['BUILD_ENV'],
    }
    _gocmd_exec = _env['GOCMD_EXEC']
    _env['FUNC_SHELL_DEVNUL'](env=_build_env, args=[_gocmd_exec, 'env'])


    _go_module = _env['FUNC_SHELL_STDOUT'](cwd=_env['PROJ_ROOT'], args=[_gocmd_exec, 'list', '-m'])[:-1]
    _go_exe = _env['FUNC_SHELL_STDOUT'](env=_build_env, args=[_gocmd_exec, 'env', 'GOEXE'])[:-1]
    _go_extra_args = _env['EXTRA_ARGS']
    _go_ldflags = "-v"

    _git_hash = _env['FUNC_SHELL_STDOUT'](
        cwd=_env['PROJ_ROOT'], args=['git', 'describe', '--tags', '--always', '--dirty', '--abbrev=7']
    )
    _go_ldflags = f"{_go_ldflags} -X '{_go_module}/system.version={_git_hash[:-1]}'"

    _build_datetime = dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%d %H:%M:%S %Z")
    _go_ldflags = f"{_go_ldflags} -X '{_go_module}/system.datetime={_build_datetime}'"

    if _env['LIB_RELEASE'] == '0':
        _go_extra_args.extend(['-gcflags', '-N -l'])
    if _env['LIB_RELEASE'] == '1':
        _go_ldflags = f"{_go_ldflags} -s -w"
    _go_ldflags = f"{_go_ldflags} -X '{_go_module}/system.flavor={_env['LIB_RELEASE']}'"


    args = [
        _gocmd_exec, 'build', #'-x',
        '-o', f"{_env['PKG_INST_DIR']}/{_go_module}{_go_exe}",
        '-ldflags', _go_ldflags,
    ]
    args.extend(_go_extra_args)
    args.append(_env['PROJ_ROOT'])
    _env['FUNC_SHELL_DEVNUL'](cwd=_env['PROJ_ROOT'], env=_build_env, args=args)
