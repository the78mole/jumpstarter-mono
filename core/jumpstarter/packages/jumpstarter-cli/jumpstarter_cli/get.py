import click
from jumpstarter_cli_common.config import opt_config
from jumpstarter_cli_common.exceptions import handle_exceptions_with_reauthentication
from jumpstarter_cli_common.opt import OutputType, opt_output_all
from jumpstarter_cli_common.print import model_print

from .common import opt_selector
from .login import relogin_client


@click.group()
def get():
    """
    Display one or many resources
    """


@get.command(name="exporters")
@opt_config(exporter=False)
@opt_selector
@opt_output_all
@click.option("--with", "with_options", multiple=True, help="Include additional information (e.g., 'leases')")
@handle_exceptions_with_reauthentication(relogin_client)
def get_exporters(config, selector: str | None, output: OutputType, with_options: tuple[str, ...]):
    """
    Display one or many exporters
    """

    include_leases = "leases" in with_options
    exporters = config.list_exporters(filter=selector, include_leases=include_leases)

    model_print(exporters, output)


@get.command(name="leases")
@opt_config(exporter=False)
@opt_selector
@opt_output_all
@handle_exceptions_with_reauthentication(relogin_client)
def get_leases(config, selector: str | None, output: OutputType):
    """
    Display one or many leases
    """

    leases = config.list_leases(filter=selector)

    model_print(leases, output)
