import 'zone.js/node';
import { existsSync } from 'fs';
import { join } from 'path';

const DIST_FOLDER = join(process.cwd(), 'dist/viewer-app/browser');

export const angularServerConfig = {
  views: DIST_FOLDER
};
