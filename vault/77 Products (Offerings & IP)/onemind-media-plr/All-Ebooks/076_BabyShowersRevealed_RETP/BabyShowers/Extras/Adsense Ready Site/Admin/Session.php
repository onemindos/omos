<?php
	
	session_start();
	
	if (isset($_SESSION['UserName']))
		{
			if (isset($_GET['Status']))
				{
					session_destroy();
					header('location: index.php');
					die();
				}
		}
	else
		{
			session_destroy();
			header('location: index.php');
			die();
		}
?>