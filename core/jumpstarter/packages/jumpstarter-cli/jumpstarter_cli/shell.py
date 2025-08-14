import sys
from datetime import timedelta

import click
from jumpstarter_core_cli_common.config import opt_config
from jumpstarter_core_cli_common.exceptions import handle_exceptions_with_reauthentication

from .common import opt_duration_partial, opt_selector
from .login import relogin_client
from jumpstarter_core.common.utils import launch_shell
from jumpstarter_core.config.client import ClientConfigV1Alpha1
from jumpstarter_core.config.exporter import ExporterConfigV1Alpha1


@click.command("shell")
@opt_config()
@click.argument("command", nargs=-1)
# client specific
# TODO: warn if these are specified with exporter config
@click.option("--lease", "lease_name")
@opt_selector
@opt_duration_partial(default=timedelta(minutes=30), show_default="00:30:00")
@click.option("--exporter-logs", is_flag=True, help="Enable exporter log streaming")
# end client specific
@handle_exceptions_with_reauthentication(relogin_client)
def shell(config, command: tuple[str, ...], lease_name, selector, duration, exporter_logs):
    """
    Spawns a shell (or custom command) connecting to a local or remote exporter

    COMMAND is the custom command to run instead of shell.

    Example:

    .. code-block:: bash

        $ jmp shell --exporter foo -- python bar.py
    """

    match config:
        case ClientConfigV1Alpha1():
            exit_code = 0
            def _launch_remote_shell(path: str) -> int:
                return launch_shell(
                    path,
                    "remote",
                    config.drivers.allow,
                    config.drivers.unsafe,
                    config.shell.use_profiles,
                    command=command,
                )

            with config.lease(selector=selector, lease_name=lease_name, duration=duration) as lease:
                with lease.serve_unix() as path:
                    with lease.monitor():
                        if exporter_logs:
                            with lease.connect() as client:
                                with client.log_stream():
                                    exit_code = _launch_remote_shell(path)
                        else:
                            exit_code = _launch_remote_shell(path)
            # we exit here to make sure that all the with clauses unwind
            sys.exit(exit_code)

        case ExporterConfigV1Alpha1():
            with config.serve_unix() as path:
                # SAFETY: the exporter config is local thus considered trusted
                launch_shell(
                    path,
                    "local",
                    allow=[],
                    unsafe=True,
                    use_profiles=False,
                    command=command,
                )
