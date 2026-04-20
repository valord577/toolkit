#!/usr/bin/env python3

# fmt: off

import sys
sys.dont_write_bytecode = True

from scripts import utils as x
# ----------------------------


import datetime as dt
import os
import shutil

from pathlib import Path
from typing import NoReturn, Union

if sys.version_info < (3, 6):
    raise RuntimeError(f'Required Python Interpreter ≥ 3.6')



# ----------------------------
# optimize
#  - 0 DEBUG
#  - 1 RELEASE (default)
# ----------------------------
LIB_RELEASE = os.getenv('LIB_RELEASE') or '1'
if LIB_RELEASE != '0': LIB_RELEASE = '1'
# ----------------------------


class _ctx:
    def __init__(self):
        self.gocmd_exec  = shutil.which('go') or 'go'
        self.extra_args  = []
        self.goenv_buld  = {}
        self.go_ldflags  = '-v'
        self.optimization = 'default'
        if LIB_RELEASE == '1':
            self.go_ldflags += ' -s -w'
            self.optimization = 'release'
        if LIB_RELEASE == '0':
            self.extra_args.extend(['-gcflags', '-N -l'])
            self.optimization = 'debug'

        self.target_plat = ''
        self.target_arch = ''
        self.target_libc = ''

        if not x.ON_GITHUB_CI:
            self.goenv_buld.update({
                'GO111MODULE': 'on',
                'GOSUMDB': 'sum.golang.google.cn',
                'GOPROXY': 'https://goproxy.cn,direct',
            })

    def getenv(self) -> dict:
        _target_archlibc = self.target_arch
        if self.target_libc:
            _target_archlibc = f'{self.target_arch}-{self.target_libc}'
        _pkg_inst_dir = (Path(x.PROJ_ROOT) / 'out' / self.target_plat / _target_archlibc).absolute().as_posix()
        if (x.ON_GITLAB_CI or x.ON_GITHUB_CI) and (_pkg_inst_dir_ci := os.getenv('INST_DIR')):
            _pkg_inst_dir = Path(_pkg_inst_dir_ci).absolute().as_posix()
        if x.ON_CODE_EDIT:
            _pkg_inst_dir = (Path(x.PROJ_ROOT) / 'out').absolute().as_posix()

        return {
            **{

                'GOCMD_EXEC': self.gocmd_exec,
                'EXTRA_ARGS': self.extra_args,
                'GOENV_BULD': self.goenv_buld,
                'GO_LDFLAGS': self.go_ldflags,
                'OPTIMIZATION': self.optimization,

                'PKG_PLATFORM': self.target_plat,
                'PKG_ARCH': self.target_arch,
                'PKG_LIBC': self.target_libc,
                'PKG_ARCH_LIBC': self.target_arch,

                'PKG_INST_DIR': _pkg_inst_dir,
            },
        }


def _setctx_linux(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    if _native:
        ctx.target_arch = x.NATIVE_ARCH
        if not (ctx.target_arch in ['arm64', 'amd64']):
            raise NotImplementedError(f'unsupported target arch: {ctx.target_arch}')
    else:
        CROSS_TOOLCHAIN_ROOT = x._util_get_cross_toolchain_dir()

        ctx.target_arch = _tuple[2]
        ctx.target_libc = _tuple[3]
        _target_triple = {
            'arm64': f'aarch64-unknown-linux-{ctx.target_libc}',
            'amd64': f'x86_64-pc-linux-{ctx.target_libc}',
            'armv7': f'arm-unknown-linux-{ctx.target_libc}',
        }[ctx.target_arch]

        # cgotool bin
        CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.getenv('CROSS_TOOLCHAIN_CGOTOOL_PREFIX')
        if not CROSS_TOOLCHAIN_CGOTOOL_PREFIX:
            CROSS_TOOLCHAIN_CGOTOOL_PREFIX = (Path(CROSS_TOOLCHAIN_ROOT) / 'cgotool-wrapper').absolute().as_posix()
        ctx.gocmd_exec = f'{CROSS_TOOLCHAIN_CGOTOOL_PREFIX}.{_target_triple}'
def _setctx_apple(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    if _native:
        ctx.target_arch = x.NATIVE_PLAT
        if not (ctx.target_arch in ['arm64', 'amd64']):
            raise NotImplementedError(f'unsupported target arch: {ctx.target_arch}')
    else:
        ctx.target_arch = _tuple[1]

        crossfiles_dir = (Path(x.PROJ_ROOT) / '.crossfiles' / 'apple')
        # cgotool bin
        ctx.gocmd_exec = (Path(crossfiles_dir) / f'cgotool-wrapper.{ctx.target_arch}').absolute().as_posix()
def _setctx_win32_mingw(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    ctx.goenv_buld.update({
        'GOOS': 'windows',
    })

    CROSS_TOOLCHAIN_ROOT = x._util_get_cross_toolchain_dir()

    ctx.target_arch = _tuple[1]
    _target_arch = {
        'arm64': f'aarch64',
        'amd64': f'x86_64',
    }[ctx.target_arch]

    # cgotool bin
    CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.getenv('CROSS_TOOLCHAIN_CGOTOOL_PREFIX')
    if not CROSS_TOOLCHAIN_CGOTOOL_PREFIX:
        CROSS_TOOLCHAIN_CGOTOOL_PREFIX = (Path(CROSS_TOOLCHAIN_ROOT) / 'cgotool-wrapper').absolute().as_posix()
    ctx.gocmd_exec = f'{CROSS_TOOLCHAIN_CGOTOOL_PREFIX}.{_target_arch}'


_targets = {
    'linux': {
        'native': True,
        'hostos': ('linux', ),
        'setctx': _setctx_linux,
        'tuples': [
            ('linux', 'crossbuild', 'amd64', 'gnu'),
            ('linux', 'crossbuild', 'arm64', 'gnu'),
            ('linux', 'crossbuild', 'armv7', 'gnueabihf'),
            ('linux', 'crossbuild', 'amd64', 'musl'),
            ('linux', 'crossbuild', 'arm64', 'musl'),
            ('linux', 'crossbuild', 'armv7', 'musleabihf'),
        ],
    },
    'darwin': {
        'native': True,
        'hostos': ('darwin', ),
        'setctx': _setctx_apple,
        'tuples': [
            ('darwin', 'arm64'),
            ('darwin', 'amd64'),
        ],
    },
    'windows': {
        'native': False,
        'hostos': ('linux', ),
        'setctx': _setctx_win32_mingw,
        'tuples': [
            ('windows', 'arm64'),
            ('windows', 'amd64'),
        ],
    },
}

def show_help(exitcode = 1) -> NoReturn:
    _native_flag_width = 0
    for k, v in _targets.items():
        _width = len(k) + 1
        if v['native'] and (_width > _native_flag_width):
            _native_flag_width = _width

    _targets_help_str = ''
    for k, v in _targets.items():
        _targets_help_str += f'    {k.ljust(_native_flag_width)}{"(* native)" if v["native"] else ""}\n'
        for tgt in v['tuples']:
            _targets_help_str += f'        {" ".join(tgt[1:])}\n'

    help_str  = f'Usage: {sys.argv[0]} -h|--help\n'
    help_str += f'Usage: {sys.argv[0]} [target]\n\n'
    help_str += f'Target Options:\n{_targets_help_str}\n'
    x.print_stderr(help_str[:-1])
    sys.exit(exitcode)


if __name__ == "__main__":
    argv_tgt: list[str] = []
    argv = sys.argv[1:]; argc = len(argv); i = 0
    while i < argc:
        arg = argv[i]; i += 1
        if arg.startswith('-h') or arg.startswith('--help'):
            show_help(0)  # exited
        else:
            argv_tgt.append(arg)

    argc_tgt = len(argv_tgt)
    if argc_tgt < 1:
        if x.NATIVE_PLAT in ['linux', 'darwin']:
            argc_tgt +=1; argv_tgt.append(x.NATIVE_PLAT)

    ctx = _ctx()
    ctx.target_plat = argv_tgt[0]
    _target = _targets.get(ctx.target_plat)
    if not _target:
        raise NotImplementedError(f'unsupported target platform: {ctx.target_plat}')
    _hostos = _target['hostos']
    if not isinstance(_hostos, tuple) \
        or \
        (len(_hostos) not in [1, 2]) \
        or \
        ((len(_hostos) == 1) and (
            _hostos[0] != x.NATIVE_PLAT
        )) \
        or \
        ((len(_hostos) == 2) and (
            _hostos[0] != x.NATIVE_PLAT or _hostos[1] != x.NATIVE_ARCH
        )):
        raise NotImplementedError(f'unsupported host os: {_hostos}')


    _tuple: Union[tuple[str, ...], None] = None
    if argc_tgt > 1:
        # check target tuple
        _tuple = tuple(argv_tgt)
        if not (_tuple in _target['tuples']):
            raise NotImplementedError(f'unsupported target tuple: {_tuple}')
    _is_native_build = ((argc_tgt == 1) and (_target['native']))
    if (not _is_native_build) and (not _tuple):
        raise NotImplementedError(f'unsupported native build: {ctx.target_plat}')
    _target['setctx'](ctx, _is_native_build, _tuple)


    build_env = ctx.getenv()
    _pkg_inst_dir = build_env['PKG_INST_DIR']

    build_steps = x._util_load_module(f'build_steps', ['module_init']).module_init(build_env)
    for func in build_steps:
        func()
    if not x.ON_CODE_EDIT:
        x._util_func__exec_python([
            (Path(x.PROJ_ROOT) / 'scripts' / 'tree.py').absolute().as_posix(), _pkg_inst_dir, '2'
        ])
    x.print_stderr(f'──── Build Done @{dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%d %H:%M:%S %Z")} ────')
