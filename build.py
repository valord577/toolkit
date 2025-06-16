#!/usr/bin/env python3

# fmt: off

import datetime as dt
import importlib.util
import os
import platform
import shutil
import subprocess as sp
import sys
from typing import NoReturn, Union


PROJ_ROOT = os.path.abspath(os.path.dirname(__file__))
# ----------------------------
# optimize
#  - 0 DEBUG
#  - 1 RELEASE (default)
# ----------------------------
LIB_RELEASE = os.getenv('LIB_RELEASE', '1')
if LIB_RELEASE != '0': LIB_RELEASE = '1'
# ----------------------------
# ci runtime
# ----------------------------
ON_GITLAB_CI = os.getenv('GITLAB_CI', '')      == 'true'
ON_GITHUB_CI = os.getenv('GITHUB_ACTIONS', '') == 'true'

class _ctx:
    def __init__(self):
        self.script = self._lazy_import()
        self.gocmd_exec  = shutil.which('go') or 'go'
        self.extra_args  = []

        self.native_plat = platform.system().lower()
        self.native_arch = platform.machine().lower()
        if self.native_arch == 'x86_64':  self.native_arch = 'amd64'
        if self.native_arch == 'aarch64': self.native_arch = 'arm64'

        self.target_plat = ''
        self.target_arch = ''
        self.target_libc = ''
        self.env_passthrough = {
            'BUILD_ENV': {},
        }

        if not ON_GITHUB_CI:
            self.env_passthrough['BUILD_ENV'].update({
                'GO111MODULE': 'on',
                'GOSUMDB': 'sum.golang.google.cn',
                'GOPROXY': 'https://goproxy.cn,direct',
            })

    def _lazy_import(self):
        name = 'build_steps.py'
        path = os.path.abspath(os.path.join(PROJ_ROOT, name))
        spec = importlib.util.spec_from_file_location('', path)
        if not spec:
            raise ModuleNotFoundError(f'missing `{name}`: "failed @importlib.util.spec_from_file_location"')
        module = importlib.util.module_from_spec(spec)
        try:
            spec.loader.exec_module(module)  # type: ignore
        except FileNotFoundError:
            raise ModuleNotFoundError(f'missing `{name}`: "no such file [{path}]"')
        if not hasattr(module, 'module_init'):
            raise ModuleNotFoundError(f'missing `{name}`: "no attr [module_init]"')
        return module

    def getenv(self) -> dict:
        env = {
            **self.env_passthrough,
            **{
                'PROJ_ROOT': PROJ_ROOT,

                'LIB_RELEASE': LIB_RELEASE,

                'FUNC_SHELL_DEVNUL': _util_func__subprocess_devnul,
                'FUNC_SHELL_STDOUT': _util_func__subprocess_stdout,

                'PKG_PLATFORM': self.target_plat,
                'PKG_ARCH': self.target_arch,
                'PKG_LIBC': self.target_libc,
                'PKG_ARCH_LIBC': self.target_arch,

                'GOCMD_EXEC': self.gocmd_exec,
                'EXTRA_ARGS': self.extra_args,
            },
        }
        if env['PKG_LIBC']:
            env['PKG_ARCH_LIBC'] = f"{env['PKG_ARCH']}-{env['PKG_LIBC']}"
        env['PKG_INST_DIR'] = os.path.abspath(os.path.join(PROJ_ROOT, 'out', env['PKG_PLATFORM'], env['PKG_ARCH_LIBC']))
        if ON_GITLAB_CI or ON_GITHUB_CI:
            env['PKG_INST_DIR'] = os.getenv('INST_DIR') or env['PKG_INST_DIR']

        return env


def _self_func__tree(basepath: str, depth: int = 0):
    # name, path, depth, is_last, is_symlink, is_dir
    _stack = [(basepath, basepath, -1, '-1', os.path.islink(basepath), os.path.isdir(basepath))]
    while _stack:
        _name, _path, _depth, _is_last, _is_symlink, _is_dir = _entry = _stack.pop()
        if _depth == -1:
            print(_name, file=sys.stderr)
        else:
            print(f"{'│   ' * _depth}{'└── ' if _is_last == '1' else '├── '}{_name}{f' -> {os.readlink(_path)}' if _is_symlink else ''}", file=sys.stderr)

        if (not _is_dir) or (depth > 0 and _depth + 1 >= depth):
            continue
        with os.scandir(_path) as it:
            entries = sorted(it, key=lambda e: e.name)
            entries_dir_first = [ d for d in entries if d.is_dir() ]
            entries_dir_first.extend([ f for f in entries if not f.is_dir() ])
            for i, entry in enumerate(reversed(entries_dir_first)):
                _stack.append(
                    (
                        entry.name,
                        entry.path,
                        _depth + 1,
                        '1' if i == 0 else '0',
                        entry.is_symlink(),
                        entry.is_dir(follow_symlinks=False),
                    )
                )
def _util_func__subprocess_stdout(args: Union[str, list[str]],
    cwd: Union[str, None] = None, env: Union[dict[str, str], None] = None, shell=False
) -> str:
    print(f'>>>> subprocess cmdline: {args}', file=sys.stderr)
    proc = sp.run(args=args, cwd=cwd, env=env, shell=shell, stdout=sp.PIPE, text=True)
    if proc.returncode != 0:
        print(f'>>>> subprocess exitcode: {proc.returncode}', file=sys.stderr)
        sys.exit(proc.returncode)
    return proc.stdout
def _util_func__subprocess_devnul(args: Union[str, list[str]],
    cwd: Union[str, None] = None, env: Union[dict[str, str], None] = None, shell=False
):
    print(f'>>>> subprocess cmdline: {args}', file=sys.stderr)
    proc = sp.run(args=args, cwd=cwd, env=env, shell=shell)
    if proc.returncode != 0:
        print(f'>>>> subprocess exitcode: {proc.returncode}', file=sys.stderr)
        sys.exit(proc.returncode)




def _setctx_linux(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    if ctx.native_plat != 'linux':
        raise NotImplementedError(f'unsupported host os: {ctx.native_plat}')
    ctx.env_passthrough['PLATFORM_LINUX'] = True

    if _native:
        ctx.target_arch = ctx.native_arch
        if not (ctx.target_arch in ['arm64', 'amd64']):
            raise NotImplementedError(f'unsupported target arch: {ctx.native_arch}')
    else:
        CROSS_TOOLCHAIN_ROOT = os.getenv('CROSS_TOOLCHAIN_ROOT')
        if not CROSS_TOOLCHAIN_ROOT:
            raise GeneratorExit('missing required env: `CROSS_TOOLCHAIN_ROOT`')

        ctx.target_arch = _tuple[2]
        ctx.target_libc = _tuple[3]

        _target_triple = ''
        if ctx.target_arch == 'arm64':
            _target_triple = f'aarch64-unknown-linux-{ctx.target_libc}'
        if ctx.target_arch == 'amd64':
            _target_triple = f'x86_64-pc-linux-{ctx.target_libc}'
        if ctx.target_arch == 'armv7':
            _target_triple = f'arm-unknown-linux-{ctx.target_libc}'

        # cgotool bin
        CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.getenv('CROSS_TOOLCHAIN_CGOTOOL_PREFIX')
        if not CROSS_TOOLCHAIN_CGOTOOL_PREFIX:
            CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.path.abspath(os.path.join(CROSS_TOOLCHAIN_ROOT, 'cgotool-wrapper'))
        ctx.gocmd_exec = f'{CROSS_TOOLCHAIN_CGOTOOL_PREFIX}.{_target_triple}'
def _setctx_apple(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    if ctx.native_plat != 'darwin':
        raise NotImplementedError(f'unsupported host os: {ctx.native_plat}')
    ctx.env_passthrough['PLATFORM_APPLE'] = True

    if _native:
        ctx.target_arch = ctx.native_arch
        if not (ctx.target_arch in ['arm64', 'amd64']):
            raise NotImplementedError(f'unsupported target arch: {ctx.native_arch}')
    else:
        ctx.target_arch = _tuple[1]

        crossfiles_dir = os.path.abspath(os.path.join(PROJ_ROOT, '.crossfiles', 'apple'))
        # cgotool bin
        ctx.gocmd_exec = os.path.abspath(os.path.join(crossfiles_dir, f'cgotool-wrapper.{ctx.target_arch}'))
def _setctx_win32_mingw(
    ctx: _ctx, _native: bool, _tuple: tuple[str, ...],
):
    if ctx.native_plat != 'linux':
        raise NotImplementedError(f'unsupported host os: {ctx.native_plat}')
    ctx.env_passthrough['PLATFORM_WIN32'] = True
    ctx.env_passthrough['BUILD_ENV'].update({
        'GOOS': 'windows',
    })

    CROSS_TOOLCHAIN_ROOT = os.getenv('CROSS_TOOLCHAIN_ROOT')
    if not CROSS_TOOLCHAIN_ROOT:
        raise GeneratorExit('missing required env: `CROSS_TOOLCHAIN_ROOT`')

    ctx.target_arch = _tuple[1]

    _target_arch = ''
    if ctx.target_arch == 'amd64': _target_arch = 'x86_64'
    if ctx.target_arch == 'arm64': _target_arch = 'aarch64'

    # cgotool bin
    CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.getenv('CROSS_TOOLCHAIN_CGOTOOL_PREFIX')
    if not CROSS_TOOLCHAIN_CGOTOOL_PREFIX:
        CROSS_TOOLCHAIN_CGOTOOL_PREFIX = os.path.abspath(os.path.join(CROSS_TOOLCHAIN_ROOT, 'cgotool-wrapper'))
    ctx.gocmd_exec = f'{CROSS_TOOLCHAIN_CGOTOOL_PREFIX}.{_target_arch}'


_targets = {
    'linux': {
        'native': True,
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
        'setctx': _setctx_apple,
        'tuples': [
            ('darwin', 'arm64'),
            ('darwin', 'amd64'),
        ],
    },
    'windows': {
        'native': False,
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
    print(help_str[:-1], file=sys.stderr)
    sys.exit(exitcode)


if __name__ == "__main__":
    if sys.version_info < (3, 6):
        raise GeneratorExit(f'Required Python Interpreter ≥ 3.6')


    argv_tgt: list[str] = []
    argv = sys.argv[1:]; argc = len(argv); i = 0
    while i < argc:
        arg = argv[i]; i += 1
        if arg.startswith('-h') or arg.startswith('--help'):
            show_help(0)  # exited
        else:
            argv_tgt.append(arg)
    argc_tgt = len(argv_tgt)


    ctx = _ctx()
    if argc_tgt < 1:
        if not (ctx.native_plat in ['linux', 'darwin']):
            raise NotImplementedError(f'unsupported native platform: {ctx.native_plat}')
        argc_tgt +=1; argv_tgt.append(ctx.native_plat)

    ctx.target_plat = argv_tgt[0]
    _target = _targets.get(ctx.target_plat)
    if not _target:
        raise NotImplementedError(f'unsupported target platform: {ctx.target_plat}')

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
    build_steps = ctx.script.module_init(build_env)
    for func in build_steps:
        func()
    _self_func__tree(build_env['PKG_INST_DIR'], depth=3)
    print(f'──── Build Done @{dt.datetime.now(dt.timezone.utc).strftime("%Y-%m-%d %H:%M:%S %Z")} ────', file=sys.stderr)
