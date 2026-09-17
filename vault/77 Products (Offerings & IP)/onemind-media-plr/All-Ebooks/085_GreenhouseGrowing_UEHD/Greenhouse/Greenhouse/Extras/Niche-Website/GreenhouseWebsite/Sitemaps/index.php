<?php include("../Includes/Start.php"); ?>
<html>
	<head>
		<title><?=$Category;?> Sitemap</title>
		<meta http-equiv="Content-Type"	content="text/html; charset=windows-1252">
		<meta name="Creator"			content="Niche Mania - http://www.niche-maniacs.com">
		<meta name="Copyright"			content="<?=$CopyRight;?>">
		<meta name="Description"		content="<?=$Category;?> Sitemap">
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
				<td class="tdHeader">
					<h1>
						<?=$MainTitle;?>
					</h1>
					<h3>
						<?=$SubTitle;?>
					</h3>
				</td>
			</tr>
			<tr>
				<td class="tdRow">
					<?php $Menu = "SiteMap"; include("../Includes/Menu.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdContent" style="border-right-width: 0; width: 780;">
					<br>
					<p>
						<b>
							<?=$Category;?> Sitemap
						</b>
					</p>
					<?php include("../Includes/Sitemap.php"); ?>
				</td>
			</tr>
			<tr>
				<td class="tdRow">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <a href="../Legal/index.php"><?=$Category;?> Legal Information</a>
				</td>
			</tr>
		</table>
		<?php include("../Includes/Footer.php"); ?>
		<?php include("../Includes/AdTracker.php"); ?>
	</body>
</html>