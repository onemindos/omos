<?php

	// -------------------------------------------------------------
 	// © 2000 - 2006, Niche Mania - All Rights Reserved Worldwide...
 	// -------------------------------------------------------------

	error_reporting(0);	
	
	if($Page != "Legal")
		{
			$strLegal = "<a href='../Legal/index.php'>".$Category." Legal Information</a>";
		}
	else
		{
			$strLegal = "<b>".$Category." Legal Information</b>";
		}

?>
			<tr>
				<td class="tdRow" colspan="<?=$Columns;?>">
					&copy; <?=date("Y");?>, <a href="<?=$Domain;?>"><?=$SiteName;?></a> - All Rights Reserved Worldwide | <?=$strLegal;?>
				</td>
			</tr>
		</table>
		<?php include("../Includes/Ping.php");?>
	</body>
</html>
<!-- P A S T E   E X T R A   J A V A S C R I P T   C O D E   B E L O W   T H I S   L I N E -->
<?php include("../Includes/AdTracker.php"); ?>