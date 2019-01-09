# GloFlow media management/publishing system
# Copyright (C) 2019 Ivan Trajkovic
#
# This program is free software; you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation; either version 2 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program; if not, write to the Free Software
# Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA

import os,sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import importlib
import argparse

from colored import fg,bg,attr
import pymongo
import delegator
#-------------------------------------------------------------
def main(p_log_fun):

	args_map = parse_args()

	mongodb_host_str      = args_map['mongodb_host']
	plots_dir_str         = args_map['plots_dir']
	py_stats_dirs_str     = args_map['py_stats_dirs']
	py_stats_dirs_lst     = py_stats_dirs_str.split(',')
	

	#---------------
	if args_map['run'] == 'crontab_full_init':

		cli_stats_path_str    = args_map['cli_stats_path']
		crontab_file_path_str = args_map['crontab_file_path']
		crontab__build_config(py_stats_dirs_lst,
			plots_dir_str,
			cli_stats_path_str,
			crontab_file_path_str,
			mongodb_host_str,
			p_log_fun)

		crontab__run(crontab_file_path_str,p_log_fun)
	#---------------
	elif args_map['run'] == 'crontab_config':

		cli_stats_path_str    = args_map['cli_stats_path']
		crontab_file_path_str = args_map['crontab_file_path']
		crontab__build_config(py_stats_dirs_lst,
			plots_dir_str,
			cli_stats_path_str,
			crontab_file_path_str,
			mongodb_host_str,
			p_log_fun)
	#---------------
	elif args_map['run'] == 'run_py_stat':

		
		mongo_client = get_mongodb_client(mongodb_host_str,p_log_fun)

		py_stat_str = args_map['py_stat']

		run_py_stat(py_stat_str,
			py_stats_dirs_lst,
			plots_dir_str,
			mongo_client,
			p_log_fun)
	#---------------

#----------------------------------------------
def run_py_stat(p_py_stat_str,
			p_py_stats_dirs_lst,
			p_plots_dir_str,
			p_mongo_client,
			p_log_fun):
	p_log_fun('FUN_ENTER','cli_stats.run_py_stat()')
	assert isinstance(p_py_stats_dirs_lst,list)

	m = import_module(p_py_stat_str,p_py_stats_dirs_lst,p_log_fun)

	img_filepath_str = '%s/%s.png'%(p_plots_dir_str,p_py_stat_str)
	m.run(p_mongo_client,p_log_fun,p_output_img_str=img_filepath_str)
#----------------------------------------------
def get_mongodb_client(p_host_str,
				p_log_fun):
	p_log_fun('FUN_ENTER','cli__gf_crawl__stats.get_mongodb_client()')

	mongo_client = pymongo.MongoClient(p_host_str,27017)
	return mongo_client
#----------------------------------------------
def import_module(p_py_stat_str,
			p_py_stats_dirs_lst,
			p_log_fun):
	assert isinstance(p_py_stats_dirs_lst,list)

	for d in p_py_stats_dirs_lst:
		abs_d = os.path.abspath(d)
		if not abs_d in sys.path: sys.path.append(abs_d)

	m = importlib.import_module(p_py_stat_str)
	return m

	run_frequency_minute_int = int(m.freq().strip().replace('m',''))
#----------------------------------------------
def crontab__build_config(p_py_stats_dirs_lst,
					p_plots_dir_str,
					p_cli_stats_path_str,
					p_crontab_file_path_str,
					p_mongodb_host_str,
					p_log_fun):
	p_log_fun('FUN_ENTER','cli_stats.crontab__build_config()')
	assert isinstance(p_py_stats_dirs_lst,list)
	for d in p_py_stats_dirs_lst: assert os.path.isdir(d)
	assert os.path.isdir(p_plots_dir_str)
	assert os.path.isfile(p_cli_stats_path_str)
	
	#-------------------------
	#GET_PY_STAT_FILES

	py_stat__files_lst = []

	for dir_str in p_py_stats_dirs_lst:

		abs_dir_str = os.path.abspath(dir_str)
		for stat_f_str in os.listdir(abs_dir_str):
			if stat_f_str.endswith('.py'):

				py_stat__file_name_str        = os.path.basename(stat_f_str)
				py_stat__file_name_no_ext_str = py_stat__file_name_str.split('.')[0]
				m                             = import_module(py_stat__file_name_no_ext_str,p_py_stats_dirs_lst,p_log_fun)
				py_stat__files_lst.append((m,py_stat__file_name_str))
	#-------------------------

	f = open(p_crontab_file_path_str,'w+') #create file
	f.write('''

#AUTOGENERATED!! - dont change manually, will be overwritten next time the generator is run

# turn off emailing on job execution
MAILTO=""
#-----------------------------------------------
#IMPORTANT!! - /dev/stdout - when run in a container this device is stdout of the container, 
#                            so that logs can be viewed prooperly with "docker logs" or "kubectl logs"

#-env_var_args=true - because these cronjobs are run in a container, where ENV vars are defined and should be 
#                     parsed by cli__gf_crawl__stats.py
#-----------------------------------------------
''')
	#-------------------------

	py_stats__names_lst = []
	for (m,py_f_str) in py_stat__files_lst:

		print ''
		print fg('green')+py_f_str+attr(0)
		print ''

		#-------------------------
		#CRON
		py_stat__file_name_no_ext_str = os.path.basename(py_f_str).split('.')[0]
		py_stat__name_str             = py_stat__file_name_no_ext_str
		
		py_stats__names_lst.append(py_stat__name_str)

		run_frequency_minute_int = int(m.freq().strip().replace('m',''))

		py_c_str = [
			'python',
			os.path.abspath(p_cli_stats_path_str),
			'-run=run_py_stat',
			'-py_stat=%s'%(py_stat__name_str),
			'-py_stats_dirs=%s'%(','.join(p_py_stats_dirs_lst)),
			'-plots_dir=%s'%(p_plots_dir_str),
			'-mongodb_host=%s'%(p_mongodb_host_str),
			'|',
			'tee /home/gf/logs/log__%s.log /dev/stdout'%(py_stat__name_str)
		]

		#minute   hour   day   month   dayofweek   command"
		#
		#minute    - any integer from 0 to 59
		#hour      - any integer from 0 to 23
		#day       - any integer from 1 to 31 (must be a valid day if a month is specified) 
		#month     - any integer from 1 to 12 (or the short name of the month such as jan or feb) 
		#dayofweek - any integer from 0 to 7, where 0 or 7 represents Sunday (or the short name of the week such as sun or mon) 
		#"/" - used to specify step values
		#      in the minute field "/n" means run command every n-th minute

		if run_frequency_minute_int == 0:

			#IMPORTANT!! - for once an hour, CRON has a different notation
			min_str = str(run_frequency_minute_int)
		else:
			min_str = '*/%d'%(run_frequency_minute_int)

		cron_line_str = '%s * * * * %s\n\n\n'%(min_str,' '.join(py_c_str))
		print 'cron_line_str - %s%s%s'%(fg('yellow'),cron_line_str,attr(0))
		#-------------------------

		f.write(cron_line_str)
	f.close()

	#-------------------------
	#IMPORTANT!! - set the generated crontab file for usage by crond
	r = delegator.run('crontab %s'%(p_crontab_file_path_str))
	if not r.out == "": print r.out
	if not r.err == "": print r.err
	#-------------------------

	return py_stats__names_lst
#----------------------------------------------
def crontab__run(p_crontab_file_path_str,
			p_log_fun):
	p_log_fun('FUN_ENTER','cli_stats.crontab__run()')
	assert isinstance(p_crontab_file_path_str,basestring)
	assert os.path.isfile(p_crontab_file_path_str)

	print ''
	print '   STARTING CROND ---------------'
	print ''

	c_str = 'crond -f'
	print c_str

	r = delegator.run(c_str)
	if not r.out == "": print r.out
	if not r.err == "": print r.err

#----------------------------------------------
#ADD!! - figure out a smarter way to pick the right hostport from p_host_port_lst,
#        instead of just picking the first element

def get_mongodb_client(p_host_str,
				p_log_fun):
	p_log_fun('FUN_ENTER','cli_stats.get_mongodb_client()')

	mongo_client = pymongo.MongoClient(p_host_str,27017)
	return mongo_client
#-------------------------------------------------------------
def parse_args():

	arg_parser = argparse.ArgumentParser(formatter_class = argparse.RawTextHelpFormatter)
	#---------------------------------
	arg_parser.add_argument('-run', 
					action  = "store",
					default = None,
					help    = '''
- %scrontab_full_init%s - build out a crontab config file and run cron itself
- %scrontab_config%s    - build out a crontab config file only
- %srun_py_stat%s       - run specified py_stat (.py script file)
					'''%(fg('yellow'),attr(0),
					fg('yellow'),attr(0),
					fg('yellow'),attr(0)))
	#---------------------------------
	arg_parser.add_argument('-mongodb_host',
					action  = "store",
					default = '127.0.0.1',
					help    = '''
host of the Mongodb server
					''')
	#---------------------------------
	arg_parser.add_argument('-py_stat',
					action  = "store",
					default = None,
					help    = '''
name of the py_stat (.py script file, but without extension) to run
					''')
	#---------------------------------
	arg_parser.add_argument('-py_stats_dirs', 
					action  = "store",
					default = './stats/',
					help    = '''
dirs in which to look for py stats scripts (.py) - ',' separated list
					''')
	#---------------------------------
	arg_parser.add_argument('-plots_dir', 
					action  = "store",
					default = './plots/',
					help    = '''
dir in which to place generated plots
					''')
	#---------------------------------
	arg_parser.add_argument('-cli_stats_path',
					action  = "store",
					default = './cli_stats.py',
					help    = '''
filepath to the cli_stats.py python script (used for crontab config to run it from the right path)
					''')
	#---------------------------------
	arg_parser.add_argument('-crontab_file_path', 
					action  = "store",
					default = './crontab.txt',
					help    = '''
crontab configuration file path (cron runs py stats scripts)
					''')
	#---------------------------------

	passed_in_args_lst = sys.argv[1:]
	args_namespace     = arg_parser.parse_args(passed_in_args_lst)

	return {
		'run':              args_namespace.run,
		'mongodb_host':     args_namespace.mongodb_host,
		'py_stat':          args_namespace.py_stat,
		'py_stats_dirs':    args_namespace.py_stats_dirs,
		'plots_dir':        args_namespace.plots_dir,
		'cli_stats_path':   args_namespace.cli_stats_path,
		'crontab_file_path':args_namespace.crontab_file_path
	}
#-------------------------------------------------------------
if __name__ == '__main__':
	def log_fun(g,m):print '%s:%s'%(g,m)
	main(log_fun)