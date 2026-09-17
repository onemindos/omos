<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title> Sudoku For Kids | Index of All Articles in the  Sudoku For Kids Category</title>
		<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Category"			content="">
		<meta name="Creator"			content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Index of All Articles in the  Sudoku For Kids Category">
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
				<td class="tdHeader" colspan="2">
					<h1>
						<?=$MainTitle;?>
					</h1>
					<h3>
						<?=$SubTitle;?>
					</h3>
				</td>
			</tr>
			<tr>
				<td class="tdRow" colspan="2">
					<?php $Menu = "Articles"; include("../Includes/Menu.php"); ?>
				</td>
			</tr>
			<tr>
				<td></td>
				<td class="tdContent" rowspan="2">
					<br>
					<div style="float: right;"> 
						<a href="../RssFeed.xml" target="_new" style="border: 0; padding-right: 5;">
							<img align="absmiddle" border="0" src="http://www.niche-mania.com/Images/Rss_32x32.jpg">
						</a>					
					</div>
					<p>
						<b>Index of All Articles in the  Sudoku For Kids Category</b> 
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<br>
<br>
<b>
<a href='../Articles/Sudoku_for_kids_and_beginners.php'>Sudoku for kids and beginners</a>
</b>
<br>
<br>
Even your four year old child who is eager to learn games can play sudoku for kids In the UK, the Sudoku fever has captured the minds of the people even the kids The sudoku game has changed the way other people look at how the game is played In the coming <br>
<small>
<a href='../Articles/Sudoku_for_kids_and_beginners.php'>Read More...</a>
</small>
<br>
<br>
 
					</p>
					<br>
				</td>
			</tr>
			<tr>
				<td class="tdLeft">
					<?php include("../Includes/J_Box.php"); ?> <?php 
					include("../Includes/Navigation.php"); ?> <?php 
					include("../Includes/Google_160x600.php"); ?>
				</td>
			</tr>
			<?php if($ShowNewsFeed) { ?>
			<tr>
				<td class="tdRow" colspan="2">
					<?=$Category;?> News and Events
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<?php include("../Includes/Google_Search.php"); ?>
				</td>
			</tr>
			<tr>
				<td colspan="2">
					<br>
					<?php include("../Includes/NewsFeed.php"); ?> <?php if ($DisplayAmazon) { echo 
					"<hr>"; echo "<br><center>"; include("../Includes/Amazon_728x90.php"); echo 
					"</center><br>"; } ?>
				</td>
			</tr>
			<?php } ?>
			<tr>
				<td class="tdRow" colspan="2">
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