<?php
/**
 * PHP探针
 * 
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
defined('DS') or define('DS',DIRECTORY_SEPARATOR);
//#函数
function color($v,$color='red'){
	return "<b style=\"color:{$color};\">{$v}</b>";
}
function color_red($v){
	return color($v,'red');
}
function color_green($v){
	return color($v,'green');
}
function url_action($ac){
	return '?ac='.$ac;
}
function qini_get($key){
	$v=ini_get($key);
	if($v===null || $v===''){
		return '<i>no value</i>';
	}else{
		return $v;
	}
}
/**
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
interface TzInterface{
	/**
	 * 
	 */
	public function tz();
}
abstract  class TzBase implements TzInterface{
	/**
	 * @param string $title
	 * @param array $datas
	 */
	public function table($title,array $datas){
		$body="";
		foreach($datas as $k=>$v){
			if($v===null || $v===''){$v='<i>no value</i>';}
			$body.="<tr><td>{$k}</td><td>{$v}</td></tr>";
		}
		$html="<table class=\"table table-hover\">
			   <thead>
				  <tr>
					<th colspan=\"2\">{$title}</th>
				  </tr>
			   </thead>
			   <tbody>
			   {$body}
			   </tbody>
			   </table>
				";
		return $html;	   
	}
	/**
	 * action="'.url_action($action).'" 
	 * 
	 * @param string $action
	 * @param string $def
	 */
	public function form($action,$def=''){
		return '<form target="_blank">
					<input type="hidden" name="ac" value="'.$action.'"/>
					<input type="text" 	 name="value" placeholder="输入..." value="'.$def.'"/>
					<button type="submit">提交</button>
				</form>
				';
	}
}
/**
 * 服务器信息
 * 
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
class TzServerInfo extends TzBase{
	/**
	 * 
	 */
	public function tz(){
		$info=[];
		$info['主机和端口']	=$_SERVER['SERVER_NAME'].':'.$_SERVER['SERVER_PORT'].' / '.$_SERVER['SERVER_ADDR'].':'.$_SERVER['SERVER_PORT'];
		$info['主机名']	=gethostname();
		$info['探针文件']	=$_SERVER['SCRIPT_FILENAME'];
		$info['标识']		=@php_uname();
		$info['操作系统']	=PHP_OS;
		$info['解译引擎']	=$_SERVER['SERVER_SOFTWARE'];
		$info['时间和时区']	=date('Y-m-d H:i:s').' 默认时区:'.date_default_timezone_get();
		//var_dump($info);
		//var_dump($_SERVER);
		echo $this->table('服务器信息',$info);
	}
}
/**
 * 服务器信息
 * 
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
class TzPhpInfo extends TzBase{
	/**
	 * get_cfg_var
	 * ini_get
	 * ini_set('memory_limit',0);
	 * var_dump(ini_get('memory_limit'));
	 * var_dump(get_cfg_var('memory_limit'));
	 */
	public function tz(){
		$info=[];
		$info['PHP版本/PHP_VERSION']			=PHP_VERSION;
		$info[color_red('运行模式/SAPI')]		=php_sapi_name();
		$info['脚本占用最大内存/memory_limit']		=ini_get('memory_limit');
		$info['POST最大数据长度/post_max_size']		=ini_get('post_max_size');
		$info['上传文件最大限制/upload_max_filesize']	=ini_get('upload_max_filesize');
		$info['浮点型数据显示的有效位数/precision']		=ini_get('precision');
		$info['脚本超时时间/max_execution_time']	=ini_get('max_execution_time');
		$info['socket超时时间/default_socket_timeout']	=ini_get('default_socket_timeout');
		
		$info['PHP页面根目录/doc_root']			=ini_get('doc_root');
		$info['用户根目录/user_dir']			=ini_get('user_dir');
		$info['显示错误信息/display_errors']	=ini_get('display_errors');
		$info['注册全局变量/register_globals']	=ini_get('register_globals');
		
		$info['短标签模式"&lt;?...?&gt;"/short_open_tag']	=ini_get('short_open_tag');
		$info['报告内存泄漏/report_memleaks']				=ini_get('report_memleaks');
		$info['允许打开远程文件/allow_url_fopen']				=ini_get('allow_url_fopen');
		$info['注册命令行变量argv,argc/register_argc_argv']	=ini_get('register_argc_argv');
		/*
		$info['']	=ini_get('');
		*/
		//var_dump($info);
		echo $this->table('PHP信息',$info);
		//
		$this->safe();
		$this->perf();
		$this->phpinfo_imp();
		$this->extensions();
	}
	/**
	 * php安全
	 */
	public function safe(){
		$info=[];
		$info[color_red('禁用危险函数/disable_functions')]	=ini_get('disable_functions');
		$info['dl()函数/enable_dl']							=ini_get('enable_dl');
		$info['eval语言结构']								=color_red('极度危险，不推荐在代码中使用。').'不是函数，不能用disable_functions禁止，需要使用第三方扩展Suhosin等';
		$info['PHP安全模式/safe_mode']						=ini_get('safe_mode');
		$info['允许访问的文件路径/open_basedir']			=ini_get('open_basedir');
		$info['']											='文件系统安全：限制访问超出open_basedir定义的路径文件，避免敏感文件被访问，passwd/php配置文件/my.ini等';
		//危险函数
		$funcs='system,show_source,shell_exec,exec,proc_open,popen,passthru,dl,assert';
		$html='';
		foreach(explode(',',$funcs) as $func){
			$f=$func;
			if(function_exists($func)){
				$func=color_green($func.'-存在');
			}else{
				$func=color_red($func.'-不存在');
			}
			$html.=$func;
		}
		$info[color_red('危险函数是否存在/启用')]="<div class=\"danger-funcs\">{$html}</div>";
		
		echo $this->table('PHP安全',$info);
	}
	/**
	 * php性能
	 * 
	 * @link http://php.net/manual/zh/opcache.configuration.php
	 */
	public function perf(){
		$opcache="查看phpinfo()信息的opcache扩展部分";
		$opcache.="\n<br/>是否开启：".qini_get('opcache.enable');
		$opc=[];
		$opc['opcache.enable/启用操作码缓存']=qini_get('opcache.enable');
		$opc['opcache.memory_consumption/共享内存大小']=qini_get('opcache.memory_consumption');
		$opc['opcache.validate_timestamps/定时检查脚本是否更新']=qini_get('opcache.validate_timestamps').' (单位秒)';
		$opc['opcache.save_comments/是否缓存注释内容']=qini_get('opcache.save_comments');
		$opc['opcache.blacklist_filename/不缓存的黑名单']=qini_get('opcache.blacklist_filename');
		$opc['opcache.force_restart_timeout/重启过期时间']=qini_get('opcache.force_restart_timeout');
		$opc['opcache.preferred_memory_model/首选缓存模块']=qini_get('opcache.preferred_memory_model').' (mmap，shm, posix 以及 win32)';
		$opc['opcache.file_cache_only/只用文件缓存还是开启共享内存']=qini_get('opcache.file_cache_only');
		
		$opc=$this->table('查看phpinfo()信息的opcache扩展部分',$opc);
		
		$info=[];
		$info[color_green('opcache/操作码缓存')]=$opc;
		echo $this->table('PHP性能',$info);
	}
	/**
	 * 已加载模块
	 */
	public function extensions(){
		$exts=get_loaded_extensions();
		$text='<div class="extensions">';
		foreach($exts as $key=>$value){
			$text.="<span>$value</span>";
		}
		$text.='</div>';
		echo $this->table('PHP已加载扩展/get_loaded_extensions',[''=>$text]);
	}
	/**
	 * phpinfo的重要字段信息
	 */
	public function phpinfo_imp(){
		$info=[];
		$info['源码编译器/Compiler']			='例如: MSVC14 (Visual C++ 2015)';
		$info['php位数/Architecture']			='例如: x86 x64';
		$info['运行模式/Server API']			='例如: Apache 2.0 Handler PHP-FPM FastCGI';
		$info['当前使用的php.ini/Loaded Configuration File']='例如: R:\xampp\php\php.ini';
		$info['是否开启操作码缓存/OPcache']		='例如: ';
		$info[color_green('是否是线程安全的/Thread Safety')]	='例如: enabled';
		$info[color_green('查看扩展库是否支持')]				='例如: gd mbstring pdo mysqli redis session';
		$info[color_red('执行phpinfo函数')]='<a href="'.url_action('phpinfo').'" target="_blank" style="color:red;"><b>点击执行</b></a>';
		echo $this->table('phpinfo()函数中的重要信息',$info);
	}
}
/**
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
class TzActions extends TzBase{
	/**
	 * 
	 */
	public function tz(){
		$ac=@$_GET['ac'];
		if($ac>''){
			$method='ac_'.$ac;
			if(method_exists($this,$method)){
				$this->$method();
				exit();
			}
		}
	}
	/**
	 *
	 */
	public function ac_phpinfo(){
		phpinfo();
	}
	/**
	 *
	 */
	public function ac_functions(){
		$ac=@$_GET['ac'];
		var_dump($ac);exit();
	}
	/**
	 *
	 */
	public function ac_function_check(){
		$value=@$_GET['value'];
		if(function_exists($value)){
			echo color_green('函数存在:'.$value);
		}else{
			echo color_red('函数不存在:'.$value);
		}
	}
	/**
	 *
	 */
	public function ac_class_check(){
		$value=@$_GET['value'];
		if(class_exists($value)){
			echo color_green('类存在:'.$value);
		}else{
			echo color_red('类不存在:'.$value);
		}
	}
	/**
	 *
	 */
	public function ac_ext_check(){
		$value=@$_GET['value'];
		if(extension_loaded($value)){
			echo color_green('扩展已加载:'.$value);
		}else{
			echo color_red('扩展未加载:'.$value);
		}
	}
	/**
	 * 数据库检测
	 * $dsn='mysql:host=localhost;dbname=test;port=3306';
	 */
	public function ac_db_check(){
		$dsn	=@$_GET['dsn'];
		$user	=@$_GET['user'];
		$pwd	=@$_GET['pwd'];
// 		var_dump($_GET);exit();
// 		$dsn='mysql:host=localhost;port=3306';
// 		$user='root';
// 		$pwd ='root';
		echo "dsn: {$dsn}<br/>\n";
		echo "user: {$user}<br/>\n";
		echo "pwd: {$pwd}<br/>\n";
		try{
			$conn=new \PDO($dsn,$user,$pwd);
			if($conn){
				echo color_green('数据库连接成功');
			}else{
				echo color_red('数据库连接失败');
			}
		}catch(\PDOException $e){
			echo color_red('数据库连接失败');
			echo "<br/>\n";
			echo color_red('异常信息: '.$e->getMessage());
		}
	}
}
/**
 * @author xiaowang <736523132@qq.com>
 * @copyright Copyright (c) 2013 http://qingmvc.com
 * @license http://www.apache.org/licenses/LICENSE-2.0 Apache-2.0
 */
class TzTools extends TzBase{
	/**
	 * 
	 */
	public function tz(){
		$this->db_check();
		//
		$list=[];
		$list['函数是否存在']=$this->form('function_check','printf').' mysql/mysql_query/mb_substr';
		$list['类是否存在']	 =$this->form('class_check','mysqli').' Pdo/Reflection';
		$list['扩展是否已加载']=$this->form('ext_check','redis').' redis/memcache';
		echo $this->table('检测函数/类/扩展',$list);
	}
	/**
	 *
	 */
	public function db_check(){
		$action='db_check';
		$dns ='mysql:host=localhost;port=3306';
		$user='root';
		$pwd ='root';
		$html='<form target="_blank">
					<input type="hidden" name="ac" value="'.$action.'" />
					<input type="text" 	 name="dsn" 	value="'.$dns.'" 	style="width: 360px;"/>
					用户名：<input type="text" name="user" value="'.$user.'" 	style="width: 50px;"/>
					密码：<input type="text" 	name="pwd" 	value="'.$pwd.'" 	style="width: 50px;"/>
					<button type="submit">提交</button>
				</form>
				';
		$list=[];
		$list[color_red('mysql')]=$html;
		$list['']='如果也要检测数据库: mysql:host=localhost;port=3306;dbname=test';
		echo $this->table(color_red('数据库连接测试'),$list);
	}
}
(new TzActions())->tz();
?>
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>探针 | QingMVC</title>
<style type="text/css">
body{
	font-family: "Microsoft YaHei","微软雅黑",Arial,"Micro Hei",SimSun,"宋体",Heiti,"黑体",sans-serif;
    font-size: 14px;
	
	width: 960px;
    margin: 0 auto;	
}
table {
    width: 100%;
    padding: 0;
    margin-top:28px;
    border-collapse: collapse;
    border-spacing: 0;
    box-shadow: 1px 1px 1px #CCC;
}
td,th{
    padding: 5px 10px;
    border: 1px solid #ccc;
	word-break: break-all;
}
td>table,
th>table{
    margin: -5px -10px;
	margin:0;
    box-shadow: none;
}
th{
	background: #ddd;
}
tr>td:first-child{
	background: #efefef;
}
h1{
    font-weight: 100;
    text-align: center;
    padding: 9px 0;

    margin-bottom: 10px;
    border-bottom: 1px solid #eee;
    margin: 20px 0 10px 0;
}
i{
    color: #ccc;	
}
.extensions,
.danger-funcs
{
	word-break: break-all;
}
.extensions>span,
.danger-funcs>b
{
    padding: 0 10px;
    display: inline-block;	
}
form{
	display:inline-block;
}
</style>
</head>
<body>
<h1>PHP探针</h1>
<?php
(new TzServerInfo())->tz();
(new TzPhpInfo())->tz();
(new TzTools())->tz();
?>
<br/>
<br/>
<br/>
</body>
</html>