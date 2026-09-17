<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title>Articles Concerning Identity Theft | Index of All Articles in the Articles Concerning Identity Theft Category</title>
		<meta http-equiv="Content-Type" content="text/html; charset=windows-1252">
		<meta name="Creator" content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright" content="<?=$CopyRight;?>">
		<meta name="Description" content="Index of All Articles in the Articles Concerning Identity Theft Category">
		<meta name="Keywords" content="<?=$Keywords;?>">
		<meta name="Distribution" content="global">
		<meta name="Publisher" content="<?=$Domain;?>">
		<meta name="Rating" content="General">
		<meta name="Revisit-after" content="5 days">
		<meta name="Robots" content="index,follow">
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
					<p>
						<b>Index of All Articles in the Articles Concerning Identity Theft Category</b> 
					</p>
					<?php include("../Includes/Google_336x280.php"); ?>
					<p>
						<br>
<br>
<b>
<a href='../Articles/Be_Informed_and_UptoDate_Read_Articles_Concerning_Identity_Theft.php'>Be Informed and Up-to-Date! Read Articles Concerning Identity Theft</a>
</b>
<br>
<br>
Are you aware of identity theft and the danger it can bring to your life if you become one of its victims Perhaps not yet Do not think that your innocence regarding identity thefts subtle workings is a protection from being violated Identity theft is a mal<br>
<small>
<a href='../Articles/Be_Informed_and_UptoDate_Read_Articles_Concerning_Identity_Theft.php'>Read More...</a>
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
