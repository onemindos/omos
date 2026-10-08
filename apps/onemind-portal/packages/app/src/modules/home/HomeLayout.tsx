import { Fragment, useMemo, useState } from 'react';
import { Box, Button, Typography } from '@material-ui/core';
import { Content } from '@backstage/core-components';
import { CustomHomepageGrid } from '@backstage/plugin-home';
import { HomePageLayoutProps } from '@backstage/plugin-home-react/alpha';
import { Directory } from './homeModule';

// Use the supported layout blueprint for natural page flow. Retain Backstage's
// native editor and saved layout in the optional widget view.
export function HomeLayout({ widgets, defaultConfig }: HomePageLayoutProps) {
  const [customize,setCustomize]=useState(false);
  const gridConfig=useMemo(()=>defaultConfig?.map(item=>({component:item.component,x:item.column,y:item.row,width:item.width,height:item.height,movable:item.movable,deletable:item.deletable,resizable:item.resizable})),[defaultConfig]);
  const widget=(name:string)=>widgets.find(item=>item.name===name)?.component;
  const utilities=['HomePageWorldClock','HomePageRecentlyVisited','HomePageMostVisited'];
  return <Content className="onemind-home-content"><Box className="onemind-home-layout">
    <Box className="onemind-home-toolbar">
      <Typography variant="overline">ONE MIND · HOME</Typography>
      <Button variant="outlined" color="primary" onClick={()=>setCustomize(value=>!value)}>{customize?'Back to app launcher':'Customize widgets'}</Button>
    </Box>
    {customize ? <Box><Typography variant="body2" color="textSecondary">Your saved widget layout. Use Edit to add, resize or rearrange widgets.</Typography><CustomHomepageGrid config={gridConfig} preventDuplicateWidgets>{widgets.map((item,index)=><Fragment key={item.name??index}>{item.component}</Fragment>)}</CustomHomepageGrid></Box> : <>
      <Directory fluid utilities={<Box className="onemind-home-utilities">{utilities.map(name=><Box key={name} className="onemind-home-utility">{widget(name)}</Box>)}</Box>} />
    </>}
  </Box></Content>;
}
