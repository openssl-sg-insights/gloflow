# GloFlow application and media management/publishing platform
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

import os, sys
cwd_str = os.path.abspath(os.path.dirname(__file__))

import subprocess
import delegator
import fabric.api

sys.path.append('%s/../utils'%(cwd_str))
import gf_cli_utils

#-------------------------------------------------------------
# PULL_IMAGE
def pull_remote(p_cont_image_name_str, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.pull_image()")

	fabric.api.run("sudo docker pull %s"%(p_cont_image_name_str))

#-------------------------------------------------------------
# PUSH
def push(p_image_full_name_str,
	p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_log_fun,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_os_docker.push()")

	#------------------
	# LOGIN
	login(p_dockerhub_user_str,
		p_dockerhub_pass_str,
		p_exit_on_fail_bool = p_exit_on_fail_bool,
		p_docker_sudo_bool  = p_docker_sudo_bool)
	#------------------
	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")

	cmd_lst.extend([
		"docker push",
		p_image_full_name_str
	])

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	stdout_str, _, exit_code_int = gf_cli_utils.run_cmd(c_str)
	

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	#------------------
	# DOCKER_LOGOUT
	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
	cmd_lst.append("docker logout")
	stdout_str, _, _ = gf_cli_utils.run_cmd(" ".join(cmd_lst))
	print(stdout_str)
	#------------------

#-------------------------------------------------------------
# BUILD_IMAGE
def build_image(p_image_name_str,
	p_dockerfile_path_str,
	p_log_fun,
	p_exit_on_fail_bool = False,
	p_docker_sudo_bool  = False):
	p_log_fun("FUN_ENTER", "gf_os_docker.build_image()")

	print(p_dockerfile_path_str)
	assert os.path.isfile(p_dockerfile_path_str)
	assert "Dockerfile" in os.path.basename(p_dockerfile_path_str)

	# full_image_name_str  = "%s/%s:%s"%(p_user_name_str, p_image_name_str, p_image_tag_str)
	context_dir_path_str = os.path.dirname(p_dockerfile_path_str)

	p_log_fun("INFO", "====================+++++++++++++++=====================")
	p_log_fun("INFO", "                 BUILDING DOCKER IMAGE")
	p_log_fun("INFO", "              %s"%(p_image_name_str))
	p_log_fun("INFO", "Dockerfile     - %s"%(p_dockerfile_path_str))
	p_log_fun("INFO", "image_name_str - %s"%(p_image_name_str))
	p_log_fun("INFO", "====================+++++++++++++++=====================")

	cmd_lst = []

	# RUN_WITH_SUDO - Docker deamon runs as root, and so for docker client to be able to connect to it
	#                 without any custom config the client needs to be run with "sudo".
	#                 if some config is in place to avoid this, set p_docker_sudo_bool to False.
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker build",
		"-f %s"%(p_dockerfile_path_str),
		"--tag=%s"%(p_image_name_str),
		context_dir_path_str
	])

	c_str = " ".join(cmd_lst)
	p_log_fun("INFO", " - %s"%(c_str))

	# change to the dir where the Dockerfile is located, for the 'docker'
	# tool to have the proper context
	old_cwd = os.getcwd()
	os.chdir(context_dir_path_str)
	
	#---------------------------------------------------
	def get_image_id_from_line(p_stdout_line_str):
		p_lst = p_stdout_line_str.split(' ')

		assert len(p_lst) == 3
		image_id_str = p_lst[2]

		# IMPORTANT!! - check that this is a valid 12 char Docker ID
		assert len(image_id_str) == 12
		return image_id_str

	#---------------------------------------------------

	stdout_str, _, exit_code_int = gf_cli_utils.run_cmd(c_str)

	for line_str in stdout_str:
		if line_str.startswith("Successfully built"):
			image_id_str = get_image_id_from_line(line_str)
			return image_id_str

	# IMPORTANT!! - if "go build" returns a non-zero exit code in some environments (CI) we
    #               want to fail with a non-zero exit code as well - this way other CI 
    #               programs will flag builds as failed.
	if not exit_code_int == 0:
		if p_exit_on_fail_bool:
			exit(exit_code_int)

	# change back to old dir
	os.chdir(old_cwd)

#-------------------------------------------------------------
# LOGIN
def login(p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_exit_on_fail_bool = True,
	p_docker_sudo_bool  = False):
	assert isinstance(p_dockerhub_user_str, basestring)
	assert isinstance(p_dockerhub_pass_str, basestring)

	cmd_lst = []
	if p_docker_sudo_bool:
		cmd_lst.append("sudo")
		
	cmd_lst.extend([
		"docker", "login",
		"-u", p_dockerhub_user_str,
		"--password-stdin"
	])
	print(" ".join(cmd_lst))

	process = subprocess.Popen(cmd_lst, stdin = subprocess.PIPE, stdout = subprocess.PIPE)
	process.stdin.write(p_dockerhub_pass_str) # write password on stdin of "docker login" command
	stdout_str, stderr_str = process.communicate() # wait for command completion
	print(stdout_str)
	print(stderr_str)

	if p_exit_on_fail_bool:
		if not process.returncode == 0:
			exit()

#---------------------------------------------------
# LOGIN__remote
def login__remote(p_dockerhub_user_str,
	p_dockerhub_pass_str,
	p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.login__remote()")
	assert isinstance(p_dockerhub_user_str, basestring)
	assert isinstance(p_dockerhub_pass_str, basestring)

	#---------------------------
	# UPLOAD_PASS_FILE

	pass_f_str = "tmp_file"
	f = open(pass_f_str, "w")
	f.write(p_dockerhub_pass_str)
	f.close()
	
	fabric.api.put(pass_f_str) # upload password file
	#---------------------------
	# IMPORTANT!! - specify pasword from stdin so that it doesnt show up
	#               as a part of the final command (in logs)
	fabric.api.run("cat %s | sudo docker login -u %s --password-stdin"%(pass_f_str, p_dockerhub_user_str))
	#---------------------------
	fabric.api.run("rm %s"%(pass_f_str))
	delegator.run("rm %s"%(pass_f_str)) # clean local tmp_file that holds the dockerhub password

#---------------------------------------------------
# CLEAN_STOP__CONTAINERS
def clean_stop__containers(p_cont_image_name_str, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.clean_stop__containers()")

	#--------------------
	# STOP_CURRENT_CONTAINERS
	image_ids_str = fabric.api.run("sudo docker ps -a | grep %s | awk '{print $1}'"%(p_cont_image_name_str))
	print("image_ids_str - %s"%(image_ids_str))

	if not image_ids_str == "":
		print("    >>  image already running - %s"%(p_cont_image_name_str))
		print("    >>  stoping containers    - %s"%(p_cont_image_name_str))

		for l in image_ids_str.split("\n"):
			image_id_str = l
			fabric.api.run("sudo docker stop %s"%(image_id_str)) #stop first
			fabric.api.run("sudo docker rm %s"%(image_id_str))   #remove, to not conflict with new ones
	#--------------------



#---------------------------------------------------
def install_base_docker(p_fab_api, p_log_fun):
	p_log_fun("FUN_ENTER", "gf_os_docker.install_base_docker()")

	p_fab_api.run("sudo apt-get clean")
	p_fab_api.run("sudo apt-get update")
	p_fab_api.run("sudo apt-get install -y apt-transport-https ca-certificates curl software-properties-common")
	p_fab_api.run("sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -")
	

	##FIX!! - hardcoding to "zesty" ubuntu version (17.04) because in 17.10 at the moment (dec 10 2017) there is no docker-ce package
	##        so Im hardcdoing 17.04 just for the moment so that the compatible docker-ce package is used
	##p_fab_api.run('sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"')
	#p_fab_api.run('sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu zesty stable"')

	p_fab_api.run("sudo apt-get update")
	p_fab_api.run("sudo apt-get install -y \
		apt-transport-https \
		ca-certificates \
		curl \
		gnupg-agent \
		software-properties-common")
	p_fab_api.run("curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -")
	p_fab_api.run('sudo add-apt-repository \
		"deb [arch=amd64] https://download.docker.com/linux/ubuntu \
		$(lsb_release -cs) \
		stable"')
	p_fab_api.run("sudo apt-get update")
	p_fab_api.run("sudo apt-get install -y docker-ce docker-ce-cli containerd.io")
	#p_fab_api.run('sudo apt-get install -y docker-ce')