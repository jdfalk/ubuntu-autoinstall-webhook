// filepath: /Users/jdfalk/repos/github.com/jdfalk/ubuntu-autoinstall-webhook/web/angular/autoinstall-web/src/app/components/monaco-editor/monaco-editor.component.ts
import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { editor } from 'monaco-editor';

@Component({
    selector: 'app-monaco-editor',
    template: `<div id="editor-container" style="height: 600px; width: 100%"></div>`,
    styles: [`
    #editor-container {
      border: 1px solid #ccc;
      margin-top: 10px;
    }
  `]
})
export class MonacoEditorComponent implements OnInit {
    @Input() code: string = '';
    @Input() language: string = 'plaintext';
    @Input() readOnly: boolean = false;
    @Output() codeChange = new EventEmitter<string>();

    private editor: any;

    ngOnInit(): void {
        this.initMonaco();
    }

    private initMonaco(): void {
        const container = document.getElementById('editor-container');
        if (container) {
            this.editor = editor.create(container, {
                value: this.code,
                language: this.language,
                theme: 'vs-dark',
                automaticLayout: true,
                readOnly: this.readOnly,
                minimap: {
                    enabled: true
                }
            });

            this.editor.onDidChangeModelContent(() => {
                const value = this.editor.getValue();
                this.codeChange.emit(value);
            });
        }
    }

    updateOptions(options: any): void {
        this.editor?.updateOptions(options);
    }

    setValue(value: string): void {
        this.editor?.setValue(value);
    }

    getValue(): string {
        return this.editor?.getValue() || '';
    }

    setLanguage(language: string): void {
        const model = this.editor?.getModel();
        if (model) {
            editor.setModelLanguage(model, language);
        }
    }
}
