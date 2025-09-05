; Simple test program for 6502 emulator
; Loads 'H' into A register and writes to display

.org $8000

start:
    LDA #$48    ; Load 'H' (ASCII 72)
    STA $F001   ; Write to display
    LDA #$65    ; Load 'e' (ASCII 101)  
    STA $F001   ; Write to display
    LDA #$6C    ; Load 'l' (ASCII 108)
    STA $F001   ; Write to display
    STA $F001   ; Write 'l' again
    LDA #$6F    ; Load 'o' (ASCII 111)
    STA $F001   ; Write to display
    LDA #$0A    ; Load newline
    STA $F001   ; Write to display
    BRK         ; Break (halt)

.org $FFFC
.word start     ; Reset vector points to start