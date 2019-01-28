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

from colored import fg, bg, attr
import delegator

#---------------------------------------------------
def run_cmd(p_cmd_str):
	print(p_cmd_str)
	r = delegator.run(p_cmd_str)
	if not r.out == '': print(r.out)
	if not r.err == '': print('%sFAILED%s >>>>>>>\n%s'%(fg('red'), attr(0), r.err))