# coding:utf-8
import inspect
import logging
import signal
import sys
import gevent
import locust
from locust import log,task,runners
from locust.argument_parser import parse_options
from locust.env import Environment
from locust.log import setup_logging, greenlet_exception_logger
import locust.stats as stats
from locust.stats import print_error_report, print_percentile_stats, print_stats, stats_printer, stats_history
from locust.stats import StatsCSV, StatsCSVFileWriter
from locust.user import User
from locust.user.inspectuser import get_task_ratio_dict, print_task_ratio
from locust.util.timespan import parse_timespan
from locust.exception import AuthCredentialsError
from newWeb import WebUI

version = locust.__version__


def is_user_class(item):
    """
    Check if a variable is a runnable (non-abstract) User class
    """
    return bool(
        inspect.isclass(item)
        and issubclass(item, User)
        and item.abstract == False
    )

def create_environment(user_classes, options, events=None):
    """
    Create an Environment instance from options
    """
    return Environment(
        user_classes=user_classes,
        tags=options.tags,
        exclude_tags=options.exclude_tags,
        events=events,
        host=options.host,
        reset_stats=options.reset_stats,
        step_load=options.step_load,
        stop_timeout=options.stop_timeout,
        parsed_options=options
    )

class Dummy(User):
    @task(20)
    def hello(self):
        pass

def main():
    # 改成直接使用满足boomer中的类
    user_classes = {"name":Dummy}
    user_classes = list(user_classes.values())
    
    # 解析命令行参数
    options = parse_options()

    # 设置logging
    if not options.skip_log_setup:
        if options.loglevel.upper() in ["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"]:
            setup_logging(options.loglevel, options.logfile)
        else:
            sys.stderr.write("非法参数--loglevel. 合法值: DEBUG/INFO/WARNING/ERROR/CRITICAL\n")
            sys.exit(1)
    logger = logging.getLogger(__name__)
    greenlet_exception_handler = greenlet_exception_logger(logger)
    
    logger.info("options.master_host="+options.master_host)
    if not options.master_host:
        sys.stdout.write("请提供--master-host参数，以便通知压测机\n")
        exit(0)

    logger.warning("这里是locust hazard改造的web-ui专用版本，"
                    "[--headless、--worker、--slave、--expect-slaves、--step-time、--run-time、--web-port、"
                     "--master、--master-bind-host、 --master-bind-port]命令参数无用！\n")


    try:
        import resource 
        if resource.getrlimit(resource.RLIMIT_NOFILE)[0] < 10000:
            # Increasing the limit to 10000 within a running process should work on at least MacOS.
            # It does not work on all OS:es, but we should be no worse off for trying.
            resource.setrlimit(resource.RLIMIT_NOFILE, [10000, resource.RLIM_INFINITY])
    except:
        logger.warning("System open file limit setting is not high enough for load testing, "
                       "and the OS wouldnt allow locust to increase it by itself."
                       "See https://docs.locust.io/en/stable/installation.html"
                       "#increasing-maximum-number-of-open-files-limit for more info.")

    # create locust Environment
    environment = create_environment(user_classes, options, events=locust.events)
    
    if options.show_task_ratio:
        print("\n Task ratio per User class")
        print( "-" * 80)
        print_task_ratio(user_classes)
        print("\n Total task ratio")
        print("-" * 80)
        print_task_ratio(user_classes, total=True)
        sys.exit(0)

    if options.show_task_ratio_json:
        from json import dumps
        task_data = {
            "per_class": get_task_ratio_dict(user_classes),
            "total": get_task_ratio_dict(user_classes, total=True)
        }
        print(dumps(task_data))
        sys.exit(0)

    # 只使用master模式运行
    runner = environment.create_master_runner()
    runner.state=runners.STATE_STOPPED

    if options.csv_prefix:
        stats_csv_writer = StatsCSVFileWriter(
            environment,stats.PERCENTILES_TO_REPORT,options.csv_prefix,options.stats_history_enabled
        )
    else:
        stats_csv_writer = StatsCSV(environment,stats.PERCENTILES_TO_REPORT)


    # 开启web-ui服务
    web_host = "0.0.0.0"
    protocol = "https" if options.tls_cert and options.tls_key else "http"
    logger.info("Starting web interface at %s://%s:%s" % (protocol, '127.0.0.1', options.web_port))
    try:
        web_ui = WebUI(
            environment,
            host=web_host,
            port=options.web_port,
            masterHost=options.master_host,
            auth_credentials=options.web_auth,
            tls_cert=options.tls_cert,
            tls_key=options.tls_key,
            stats_csv_writer=stats_csv_writer,
        )
    except AuthCredentialsError:
        logger.error("Credentials supplied with --web-auth should have the format: username:password")
        sys.exit(1)
    else:
        main_greenlet = web_ui.greenlet
    
    # Fire locust init event which can be used by end-users' code to run setup code that
    # need access to the Environment, Runner or WebUI
    environment.events.init.fire(environment=environment, runner=runner, web_ui=web_ui)

    stats_printer_greenlet = None
    if not options.only_summary and (options.print_stats or (options.headless and not options.worker)):
        # spawn stats printing greenlet
        stats_printer_greenlet = gevent.spawn(stats_printer(runner.stats))
        stats_printer_greenlet.link_exception(greenlet_exception_handler)

    if options.csv_prefix:
        gevent.spawn(stats_csv_writer.stats_writer).link_exception(greenlet_exception_handler)



    def shutdown():
        """
        Shut down locust by firing quitting event, printing/writing stats and exiting
        """
        logger.info("Running teardowns...")
        environment.events.quitting.fire(environment=environment, reverse=True)
        
        # determine the process exit code
        if log.unhandled_greenlet_exception:
            code = 2
        elif environment.process_exit_code is not None:
            code = environment.process_exit_code
        elif len(runner.errors) or len(runner.exceptions):
            code = options.exit_code_on_error
        else:
            code = 0
        
        logger.info("Shutting down (exit code %s), bye." % code)
        if stats_printer_greenlet is not None:
            stats_printer_greenlet.kill(block=False)
        logger.info("Cleaning up runner...")
        if runner is not None:
            runner.quit()
        
        print_stats(runner.stats, current=False)
        print_percentile_stats(runner.stats)

        print_error_report(runner.stats)
        sys.exit(code)
    
    # install SIGTERM handler
    def sig_term_handler():
        logger.info("Got SIGTERM signal")
        shutdown()
    gevent.signal_handler(signal.SIGTERM, sig_term_handler)
    
    try:
        logger.info("Starting Locust %s" % version)
        main_greenlet.join()
        shutdown()
    except KeyboardInterrupt as e:
        shutdown()

if __name__=='__main__':
    main()