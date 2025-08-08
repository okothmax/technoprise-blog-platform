import { Component, inject, OnInit, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatDialogModule, MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatDividerModule } from '@angular/material/divider';
import { AccessibilityService } from '../../services/accessibility.service';

interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText: string;
  cancelText: string;
  confirmColor?: 'primary' | 'accent' | 'warn';
}

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  imports: [
    CommonModule,
    MatDialogModule,
    MatButtonModule,
    MatIconModule,
    MatDividerModule
  ],
  templateUrl: './confirm-dialog.component.html'
})
export class ConfirmDialogComponent implements OnInit {
  data: ConfirmDialogData;

  constructor(
    @Inject(MAT_DIALOG_DATA) public dialogData: ConfirmDialogData,
    private dialogRef: MatDialogRef<ConfirmDialogComponent>,
    private accessibilityService: AccessibilityService
  ) {
    this.data = dialogData;
  }

  ngOnInit(): void {
    // Announce dialog content for screen readers
    this.accessibilityService.announce(this.data.message);
  }
}
