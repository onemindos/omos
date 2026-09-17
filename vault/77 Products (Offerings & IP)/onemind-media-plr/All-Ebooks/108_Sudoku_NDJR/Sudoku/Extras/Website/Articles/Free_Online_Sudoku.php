<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title> Free Online Sudoku | Index of All Articles in the  Free Online Sudoku Category</title>
		<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Category"			content="">
		<meta name="Creator"			content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="Index of All Articles in the  Free Online Sudoku Category">
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
						<b>Index of All Articles in the  Free Online Sudoku Category</b> 
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<br>
<br>
<b>
<a href='../Articles/Lets_all_play_free_Online_Sudoku.php'>Let’s all play free Online Sudoku</a>
</b>
<br>
<br>
Sudoku puzzles can be played anywhere whether you are at home or outside For those who don t know Sudoku puzzles but want to learn how to play the game, they can easily learn it with the many free online sudoku resources All you have to do is to visit the <br>
<small>
<a href='../Articles/Lets_all_play_free_Online_Sudoku.php'>Read More...</a>
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