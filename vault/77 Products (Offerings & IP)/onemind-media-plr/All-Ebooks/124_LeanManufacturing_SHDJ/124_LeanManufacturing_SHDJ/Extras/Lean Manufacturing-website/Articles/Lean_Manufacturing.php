<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title> Lean Manufacturing | Index of All Articles in the  Lean Manufacturing Category</title>
		<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Category"			content="">		
		<meta name="Creator"			content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Index of All Articles in the  Lean Manufacturing Category">
		<meta name="Keywords"			content="<?=$Keywords;?>">
		<meta name="Distribution"		content="global">
		<meta name="Publisher"			content="<?=$Domain;?>">
		<meta name="Rating"				content="General">
		<meta name="Revisit-after"		content="5 days">
		<meta name="Robots"				content="index,follow">
		<link href="../Includes/Styles.css" rel="stylesheet" type="text/css">
		<script src="../Includes/JavaScript.js"></script>
	</head>
	<body>
		<table align="center" cellpadding="0" cellspacing="0" class="tblMain">
			<tr>
				<td class="tdHeader" colspan="3">
					<h1>
						<?=$MainTitle;?>
					</h1>
					<h3>
						<?=$SubTitle;?>
					</h3>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="3">
					<?php $Menu = "Articles"; include("../Includes/Menu.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdLeft">
					<?php include("../Includes/J_Box.php"); ?> 
					<?php include("../Includes/Navigation.php"); ?>
				</td>
				<td class="tdContent">
					<br>
					<div style="float: right;"> 
						<a href="../RssFeed.xml" target="_new" style="border: 0; padding-right: 5;">
							<img align="absmiddle" border="0" src="http://www.niche-mania.com/Images/Rss_32x32.jpg">
						</a>					
					</div>					
					<p>
						<b>Index of All Articles in the  Lean Manufacturing Category</b> 
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<br>
<br>
<b>
<a href='../Articles/Understanding_Lean_Manufacturing.php'>Understanding Lean Manufacturing</a>
</b>
<br>
<br>
This might approach you as a big surprise, but the principles of lean manufacturing actually dates a whole lot back, specifically the era of Benjamin Franklin This does not assume that lean manufacturing principles are ancient, but the time of Benjamin Fra<br>
<small>
<a href='../Articles/Understanding_Lean_Manufacturing.php'>Read More...</a>
</small>
<br>
<br>
<br>
<br>
<b>
<a href='../Articles/The_Philosophy_of_Lean_Manufacturing.php'>The Philosophy of Lean Manufacturing</a>
</b>
<br>
<br>
In the field of management, there are various approaches that principally aim to accelerate and boost corporate firms revenues and operational efficiencies Progressive companies from around the world have been coming up every now and then different strateg<br>
<small>
<a href='../Articles/The_Philosophy_of_Lean_Manufacturing.php'>Read More...</a>
</small>
<br>
<br>
 
					</p>
					<br>
				</td>
				<td class="tdRight">
					<?php include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<?php if($ShowNewsFeed) { ?>
			<tr>
				<td class="tdRow" colspan="3">
					<?=$Category;?> News and Events
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="3">
					<br>
					<?php include("../Includes/NewsFeed.php"); ?> <?php if ($DisplayAmazon) { echo 
					"<hr>"; echo "<br><center>"; include("../Includes/Amazon_728x90.php"); echo 
					"</center><br>"; } ?>
				</td>
			</tr>
			<?php } ?>
			<tr>
				<td class="tdRow" colspan="3">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All 
					Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal 
						Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/Footer.php"); ?> <?php 
		include("../Includes/AdTracker.php"); ?>
	</body>
</html>