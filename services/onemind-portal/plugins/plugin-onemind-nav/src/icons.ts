import { ComponentType } from 'react';
import AccountTree from '@material-ui/icons/AccountTree';
import Announcement from '@material-ui/icons/Announcement';
import Archive from '@material-ui/icons/Archive';
import Assessment from '@material-ui/icons/Assessment';
import AddCircleOutline from '@material-ui/icons/AddCircleOutline';
import Build from '@material-ui/icons/Build';
import Cloud from '@material-ui/icons/Cloud';
import Description from '@material-ui/icons/Description';
import Extension from '@material-ui/icons/Extension';
import Gavel from '@material-ui/icons/Gavel';
import History from '@material-ui/icons/History';
import Home from '@material-ui/icons/Home';
import Input from '@material-ui/icons/Input';
import LibraryBooks from '@material-ui/icons/LibraryBooks';
import LocalShipping from '@material-ui/icons/LocalShipping';
import MenuBook from '@material-ui/icons/MenuBook';
import MoneyIcon from '@material-ui/icons/MonetizationOn';
import MoreHoriz from '@material-ui/icons/MoreHoriz';
import Search from '@material-ui/icons/Search';
import Timeline from '@material-ui/icons/Timeline';
import TrackChanges from '@material-ui/icons/TrackChanges';
import ViewModule from '@material-ui/icons/ViewModule';
import Visibility from '@material-ui/icons/Visibility';
import type { IconKey } from './navModel';

export const ICONS: Record<IconKey, ComponentType> = {
  search: Search,
  home: Home,
  liveOps: Timeline,
  pillar: ViewModule,
  catalog: LibraryBooks,
  api: Extension,
  graph: AccountTree,
  cloud: Cloud,
  cost: MoneyIcon,
  delivery: LocalShipping,
  insights: Assessment,
  radar: TrackChanges,
  docs: MenuBook,
  runbook: Description,
  announce: Announcement,
  adr: Gavel,
  create: AddCircleOutline,
  register: Input,
  retire: Archive,
  audit: History,
  visualizer: Visibility,
  devtools: Build,
  more: MoreHoriz,
};
