import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

from colored import fg,bg,attr

sys.path.append('%s/..'%(cwd_str))
import cli_stats
#----------------------------------------------
def log_fun(g,m):print '%s:%s'%(g,m)
#----------------------------------------------
def test():

	#---------------
	#TEST_DATA
	test_mongodb_host_str = '127.0.0.1'
	py_stats_dirs_lst     = [
		os.path.abspath('%s/../../../apps/gf_crawl_lib/py/stats'%(cwd_str)),
		os.path.abspath('%s/../../../gf_core/py/stats'%(cwd_str)),
	]
	plots_dir_str         = '%s/plots'%(cwd_str)
	cli_stats_path_str    = '%s/../cli_stats.py'%(cwd_str)
	crontab_file_path_str = '%s/crontab.txt'%(cwd_str)
	#---------------

	mongo_client = cli_stats.get_mongodb_client(test_mongodb_host_str,log_fun)

	py_stats__names_lst = cli_stats.crontab__build_config(py_stats_dirs_lst,
					plots_dir_str,
					cli_stats_path_str,
					crontab_file_path_str,
					test_mongodb_host_str,
					log_fun)

	##START CROND DEAMON
	#crontab__run(crontab_file_path_str,log_fun)

	#---------------
	#RUN_INDIVIDUAL__PY_STATS

	for py_stat_str in py_stats__names_lst:

		print ''
		print '   %sTEST%s PY_STAT - %s%s%s   >>>>>>>>>>>>>> '%(fg('yellow'),attr(0),fg('blue'),py_stat_str,attr(0))
		print ''

		cli_stats.run_py_stat(py_stat_str,
				py_stats_dirs_lst,
				plots_dir_str,
				mongo_client,
				log_fun)
	#---------------

#----------------------------------------------
test()