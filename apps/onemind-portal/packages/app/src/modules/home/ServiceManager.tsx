import { useState } from 'react';
import { Box, Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField, Typography } from '@material-ui/core';
import { AppEntry, safeUrl, useShortcuts } from './useDirectory';

export function ServiceManager(){
  const personal=useShortcuts();
  const [open,setOpen]=useState(false),[name,setName]=useState(''),[url,setUrl]=useState(''),[editing,setEditing]=useState<string>(),[error,setError]=useState(''),[busy,setBusy]=useState(false);
  const reset=()=>{setName('');setUrl('');setEditing(undefined);setError('');};
  const save=async()=>{
    const valid=safeUrl(url.trim());
    if(!valid||!name.trim()){setError('Enter a name and an HTTPS URL without embedded credentials.');return;}
    if(personal.shortcuts.some(item=>item.url===valid&&item.name!==editing)){setError('That shortcut already exists.');return;}
    if(!editing&&personal.shortcuts.length>=100){setError('You can save up to 100 personal shortcuts.');return;}
    setBusy(true);
    try{
      const id=editing??crypto.randomUUID();
      const item:AppEntry={name:id,ref:`shortcut:${id}`,title:name.trim().slice(0,100),url:valid,description:'Personal shortcut',group:'personal-shortcuts',groupTitle:'My shortcuts',personal:true,favorite:editing?personal.shortcuts.find(value=>value.name===editing)?.favorite:true};
      await personal.save(editing?personal.shortcuts.map(value=>value.name===editing?item:value):[...personal.shortcuts,item]);
      reset();setOpen(false);
    }catch{setError('Shortcut could not be saved. Try again.');}finally{setBusy(false);}
  };
  const remove=async(id:string)=>{setBusy(true);try{await personal.save(personal.shortcuts.filter(item=>item.name!==id));}catch{setError('Could not remove shortcut.');}finally{setBusy(false);}};
  return <Box mb={2}>
    <Box display="flex" flexWrap="wrap" style={{gap:8}}>
      <Button variant="contained" color="primary" disabled={!personal.loaded} onClick={()=>{reset();setOpen(true);}}>Add shortcut</Button>
      <Button variant="outlined" color="primary" href="/create/templates/default/onemind-service-registration">Add shared service</Button>
      <Button variant="text" href="/catalog-import">Register existing service</Button>
      {!!personal.shortcuts.length&&<Button onClick={()=>{reset();setOpen(true);}}>Manage shortcuts</Button>}
    </Box>
    <Dialog open={open} onClose={()=>!busy&&setOpen(false)} fullWidth maxWidth="sm" aria-labelledby="shortcut-title">
      <DialogTitle id="shortcut-title">{editing?'Edit shortcut':'Your personal shortcuts'}</DialogTitle>
      <DialogContent>
        <Typography variant="body2">Saved to your OneMind account across devices. Shared services use the catalog.</Typography>
        <TextField fullWidth label="Name" value={name} onChange={e=>setName(e.target.value)} margin="normal" inputProps={{maxLength:100}} />
        <TextField fullWidth label="HTTPS application URL" value={url} onChange={e=>setUrl(e.target.value)} margin="normal" inputProps={{maxLength:2000}} />
        {error&&<Typography role="alert" color="error">{error}</Typography>}
        {personal.shortcuts.map(item=><Box key={item.name} mt={1} display="flex" alignItems="center" style={{gap:8,flexWrap:'wrap'}}><Typography style={{flex:1}}>{item.title}</Typography><Button disabled={busy} onClick={()=>{setEditing(item.name);setName(item.title);setUrl(item.url);}}>Edit</Button><Button disabled={busy} onClick={()=>remove(item.name)}>Remove</Button></Box>)}
      </DialogContent>
      <DialogActions><Button disabled={busy} onClick={()=>setOpen(false)}>Close</Button><Button disabled={busy||!personal.loaded} color="primary" variant="contained" onClick={save}>Save shortcut</Button></DialogActions>
    </Dialog>
  </Box>;
}
